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
	docker compose down -v

build_shell:
	docker compose build

up:
	docker compose up

shell:
	docker compose run --service-ports --rm controller bash

lint:
	golangci-lint run

fix_lint:
	golangci-lint run --fix

test:
	ENV=test go test ./...

generate_docs:
	cmd/swag init -g cmd/controller/main.go --generatedTime=false --markdownFiles docs/markdown

setup_gke_kube:
	gcloud auth activate-service-account --key-file ${SERVICE_ACCOUNT_KEY}
	gcloud config set project ${GCP_PROJECT_ID}
	gcloud container clusters get-credentials ${CLUSTER_NAME} --region ${GCP_REGION}

update_gke_controller_service:
	kubectl set image deployment/uffizzi-controller -n uffizzi controller=${CONTROLLER_IMAGE}

sentry_release:
	sentry-cli releases new ${SHORT_VERSION}
	sentry-cli releases set-commits --auto ${SHORT_VERSION} --ignore-missing
	sentry-cli releases finalize ${SHORT_VERSION}
	sentry-cli releases deploys ${SHORT_VERSION} new -e ${ENV}

get_token:
	gcloud config config-helper --format="value(credential.access_token)"
