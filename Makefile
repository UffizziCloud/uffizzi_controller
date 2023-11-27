.PHONY: release
ifeq (${BRANCH},)
BRANCH = $$(git rev-parse --abbrev-ref HEAD)
endif
VERSION = $$(git rev-parse HEAD)
SHORT_VERSION = $$(git rev-parse --short=7 HEAD)

RED='\033[1;31m'
CYAN='\033[1;36m'
NO_COLOR='\033[0m'

CURRENT_VERSION := $(shell cat version)
NEXT_PATCH := $(shell docker-compose run --rm toolbox bash -c 'CURRENT_VERSION=$(CURRENT_VERSION) && semver bump patch $$CURRENT_VERSION')
NEXT_MINOR := $(shell docker-compose run --rm toolbox bash -c 'CURRENT_VERSION=$(CURRENT_VERSION) && semver bump minor $$CURRENT_VERSION')
NEXT_MAJOR := $(shell docker-compose run --rm toolbox bash -c 'CURRENT_VERSION=$(CURRENT_VERSION) && semver bump major $$CURRENT_VERSION')

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

get_token:
	gcloud config config-helper --format="value(credential.access_token)"

version:
	echo ${CURRENT_VERSION}

release_patch: export NEW_VERSION=${NEXT_PATCH}
release_patch:
	make release

release_minor: export NEW_VERSION=${NEXT_MINOR}
release_minor:
	make release

release_major: export NEW_VERSION=${NEXT_MAJOR}
release_major:
	make release

release:
	git checkout develop
	@echo "Bumping version from $(CURRENT_VERSION) to $(NEW_VERSION)"
	echo $(NEW_VERSION) > 'version'
	@echo 'Set a new chart version'
	sed 's/^\(version: \).*$$/\1$(NEW_VERSION)/' ./charts/uffizzi-controller/Chart.yaml > temp.yaml && mv temp.yaml ./charts/uffizzi-controller/Chart.yaml
	git commit -am "Change version to $(NEW_VERSION)"
	git push origin develop
	git checkout main
	@echo 'Update remote origin'
	git remote update
	git pull origin --rebase main
	git merge --no-ff --no-edit origin/develop
	git push origin main
	@echo 'Create a new tag'
	git tag uffizzi-controller-${NEW_VERSION}
	git push origin uffizzi-controller-${NEW_VERSION}
