APP_NAME=uffizzi

CONTROLLER_IMAGE_NAME="${APP_NAME}-controller"
CI_PROJECT_ID="sapient-flare-242118"
CI_REPO_URL="gcr.io/${CI_PROJECT_ID}"
CONTROLLER_IMAGE="${CI_REPO_URL}/${CONTROLLER_IMAGE_NAME}"
GCP_REPO_URL="gcr.io/${GCP_PROJECT_ID}"
GCP_CONTROLLER_IMAGE="${GCP_REPO_URL}/${CONTROLLER_IMAGE_NAME}"

APP_PREFIX=${APP_NAME}

SHORT_TARGET=$$(git rev-parse --short=7 $${COMMIT:-$$(git rev-parse $${BRANCH:-$$(git rev-parse --abbrev-ref HEAD)})})

ifeq (${BRANCH},)
BRANCH = $$(git rev-parse --abbrev-ref HEAD)
endif
VERSION = $$(git rev-parse HEAD)
SHORT_VERSION = $$(git rev-parse --short=7 HEAD)

RED='\033[1;31m'
CYAN='\033[1;36m'
NO_COLOR='\033[0m'

# Development targets

clean_world: destroy_shell build_shell up

destroy_shell:
	docker-compose down -v

build_shell:
	docker-compose build

up:
	docker-compose up

shell:
	docker-compose run --service-ports --rm controller bash

lint:
	golangci-lint run

fix_lint:
	golangci-lint run --fix

test:
	ENV=test go test ./...

generate_docs:
	cmd/swag init -g cmd/controller/main.go --generatedTime=false --markdownFiles docs/markdown

image_names:
	export CONTROLLER_IMAGE=${CONTROLLER_IMAGE}:${VERSION}

build_controller:
	docker pull "${CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG}" || true
	docker pull "${CONTROLLER_IMAGE}:latest" || true
	docker build \
		--cache-from "${CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG}" \
		--cache-from "${CONTROLLER_IMAGE}:latest" \
		--build-arg SENTRY_RELEASE=${SHORT_VERSION} \
		-t ${CONTROLLER_IMAGE}:${VERSION} \
		-t "${CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG}" \
		-t ${CONTROLLER_IMAGE}:${SHORT_VERSION} \
		-t ${CONTROLLER_IMAGE}:latest \
	  .

push_controller:
	docker push ${CONTROLLER_IMAGE}:${VERSION}
	docker push ${CONTROLLER_IMAGE}:${SHORT_VERSION}
	docker push ${CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG} || true
	docker push ${CONTROLLER_IMAGE}:latest

pull_ci_image:
	docker pull ${CONTROLLER_IMAGE}:${VERSION}
	docker pull ${CONTROLLER_IMAGE}:${SHORT_VERSION}
	docker pull ${CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG} || true
	docker pull ${CONTROLLER_IMAGE}:latest

tag_image:
	docker tag ${CONTROLLER_IMAGE}:${SHORT_VERSION} ${GCP_CONTROLLER_IMAGE}:${VERSION}
	docker tag ${CONTROLLER_IMAGE}:${SHORT_VERSION} ${GCP_CONTROLLER_IMAGE}:${SHORT_VERSION}
	docker tag ${CONTROLLER_IMAGE}:${SHORT_VERSION} ${GCP_CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG} || true
	docker tag ${CONTROLLER_IMAGE}:${SHORT_VERSION} ${GCP_CONTROLLER_IMAGE}:latest

push_gcp_controller:
	docker push ${GCP_CONTROLLER_IMAGE}:${VERSION}
	docker push ${GCP_CONTROLLER_IMAGE}:${SHORT_VERSION}
	docker push ${GCP_CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG} || true
	docker push ${GCP_CONTROLLER_IMAGE}:latest

update_gke_controller_service:
	kubectl set image deployment/uffizzi-controller -n uffizzi-controller uffizzi-controller=${GCP_CONTROLLER_IMAGE}:${SHORT_VERSION}

sentry_release:
	sentry-cli releases new ${SHORT_VERSION}
	sentry-cli releases set-commits --auto ${SHORT_VERSION} --ignore-missing
	sentry-cli releases finalize ${SHORT_VERSION}
	sentry-cli releases deploys ${SHORT_VERSION} new -e ${ENV}
