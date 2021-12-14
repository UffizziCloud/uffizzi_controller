APP_NAME=uffizzi
CI_PROJECT_ID="sapient-flare-242118"
GCP_REGION=us-central1-c

CI_REPO_URL="gcr.io/${CI_PROJECT_ID}"
GCP_REPO_URL="gcr.io/${GCP_PROJECT_ID}"

CONTROLLER_IMAGE_NAME="${APP_NAME}-controller"
CONTROLLER_IMAGE="${CI_REPO_URL}/${CONTROLLER_IMAGE_NAME}"
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

# CI targets

image_names:
	export CONTROLLER_IMAGE=${CONTROLLER_IMAGE}:${VERSION}

ci_registry_login:
	echo ${CI_SERVICE_ACCOUNT_KEY} | docker login -u _json_key --password-stdin https://gcr.io

registry_login:
	cat ${SERVICE_ACCOUNT_KEY} | docker login -u _json_key --password-stdin https://gcr.io

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

push_gcp_controller:
	docker push ${GCP_CONTROLLER_IMAGE}:${VERSION}
	docker push ${GCP_CONTROLLER_IMAGE}:${SHORT_VERSION}
	docker push ${GCP_CONTROLLER_IMAGE}:${CI_COMMIT_REF_SLUG} || true
	docker push ${GCP_CONTROLLER_IMAGE}:latest

sentry_environment:
	export SENTRY_ORG=${APP_NAME}-cloud
	export SENTRY_PROJECT=${APP_NAME}-controller

sentry_release:
	sentry-cli releases new ${SHORT_VERSION}
	sentry-cli releases set-commits --auto ${SHORT_VERSION}
	sentry-cli releases finalize ${SHORT_VERSION}
	sentry-cli releases deploys ${SHORT_VERSION} new -e ${ENV}

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

update_image: pull_ci_image tag_image registry_login push_gcp_controller

setup_gke_kube:
	gcloud auth activate-service-account --key-file ${SERVICE_ACCOUNT_KEY}
	gcloud config set project ${GCP_PROJECT_ID}
	gcloud container clusters get-credentials ${CLUSTER_NAME} --zone ${GCP_REGION}

setup_eks_kube:
	aws eks --region ${AWS_DEFAULT_REGION} update-kubeconfig --name ${CLUSTER_NAME}

aks_login:
	az login --service-principal --username ${AZURE_USERNAME} --password ${AZURE_SERVICE_PRINCIPAL_KEY} --tenant ${AZURE_TENANT}

setup_aks_kube: aks_login
	az aks get-credentials --resource-group ${CLUSTER_NAME} --name ${CLUSTER_NAME}

update_gke_controller_service: setup_gke_kube
	kubectl set image deployment/uffizzi-controller -n uffizzi-controller controller=${GCP_CONTROLLER_IMAGE}:${SHORT_VERSION}

update_eks_controller_service: setup_eks_kube
	kubectl set image deployment/uffizzi-controller -n uffizzi-controller controller=${GCP_CONTROLLER_IMAGE}:${SHORT_VERSION}

update_aks_controller_service: setup_aks_kube
	kubectl set image deployment/uffizzi-controller -n uffizzi-controller controller=${GCP_CONTROLLER_IMAGE}:${SHORT_VERSION}

generate_docs:
	cmd/swag init -g cmd/controller/main.go --generatedTime=false --markdownFiles docs/markdown
