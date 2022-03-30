APP_NAME=uffizzi

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

generate_docs:
	cmd/swag init -g cmd/controller/main.go --generatedTime=false --markdownFiles docs/markdown
