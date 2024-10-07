VERSION=0.41.3
GIT_SHORT_HASH=$(shell git rev-parse --short HEAD)
GIT_TAG_NAME="release-"$(VERSION)
IMAGE_REGISTRY=gcr.io/tse-kakeru/
IMAGE_PATH=$(IMAGE_REGISTRY)kubernetes-history-inspector
GTAG_ID="G-JJ6G0C6V06"

BUG_REPORT_URL="https://b.corp.google.com/issues/new?component=1265687&template=1747079"
DOCUMENT_URL="http://go/khi"
GCLOUD_PROJECT="kubernetes-history-inspector"
GCLOUD=gcloud --project $(GCLOUD_PROJECT)

# Top level commands for development

## Watch
# Start development server for frontend.

# For main frontend view
.PHONY=watch-web-front
watch-web-front:
	cd web && NG_APP_BACKEND_ROOT_URL="http://localhost:8080/" NG_APP_REPORT_BUG_URL=$(BUG_REPORT_URL) NG_APP_DOCUMENT_URL=$(DOCUMENT_URL) ng serve

.PHONY=watch-web-front-viewer-mode
watch-web-front-viewer-mode:
	cd web && NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER="true" NG_APP_VIEWER_MODE="true" NG_APP_BACKEND_ROOT_URL="http://localhost:8080/" NG_APP_REPORT_BUG_URL=$(BUG_REPORT_URL) NG_APP_DOCUMENT_URL=$(DOCUMENT_URL) ng serve

## Test

.PHONY=test
test: test-web test-go

# Generate the coverage report
.PHONY=coverage
coverage: coverage-go coverage-web

# Update the snapshot data used in snapshot testing in server side
.PHONY=update-snapshot
update-snapshot:
	UPDATE_SNAPSHOT=true go test -v ./...

.PHONY=lint
lint: lint-web lint-go

.PHONY=format
format: format-web format-go

### Build
# Build the frontend artifacts in dist/ folder

.PHONY=build-web
build-web: build-web-frontend

.PHONY=build-web-beta
build-web-beta: build-web-frontend-beta

.PHONY=build-web-viewer-mode
build-web-viewer-mode: build-web-frontend-viewer-mode

### Deploy

.PHONY=deploy
deploy: deploy-precondition qa deploy-container-image
	git tag -a $(GIT_TAG_NAME) -m $(GIT_TAG_NAME)
	git push origin $(GIT_TAG_NAME):refs/for/main

.PHONY=deploy-beta
deploy-beta: deploy-container-image-beta

### Initial setup

.PHONY=setup-hooks
setup-hooks:
	cp ./scripts/pre-push .git/hooks/
	cp ./scripts/pre-commit .git/hooks/
	chmod +x .git/hooks/pre-push
	chmod +x .git/hooks/pre-commit

###############################
# For internal use only       #
# TODO: remove before OSSing  #
###############################
.PHONY=deploy-khi-ro
deploy-khi-ro:
	$(GCLOUD) builds submit --config .cloudbuild/deploy-khi-ro.yaml

# Sub level commands used in the top level commands

## Build subcommands

.PHONY=build-web-frontend
build-web-frontend: ./web/**/*.ts ./web/**/*.html ./web/**/*.sass
	cd web &&NG_APP_VERSION="$(VERSION)" NG_APP_GTAG_ID="$(GTAG_ID)" NG_APP_GOOGLE_DRIVE_CLIENT_ID="" NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER="false" NG_APP_REPORT_BUG_URL=$(BUG_REPORT_URL) NG_APP_DOCUMENT_URL=$(DOCUMENT_URL) npx ng build --output-path ../dist

.PHONY=build-web-frontend-beta
build-web-frontend-beta: ./web/**/*.ts ./web/**/*.html ./web/**/*.sass
	cd web &&NG_APP_VERSION="beta-$(VERSION)@$(GIT_SHORT_HASH)" NG_APP_GTAG_ID="$(GTAG_ID)" NG_APP_GOOGLE_DRIVE_CLIENT_ID="" NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER="false" NG_APP_REPORT_BUG_URL=$(BUG_REPORT_URL) NG_APP_DOCUMENT_URL=$(DOCUMENT_URL) npx ng build --output-path ../dist --optimization=false

# TODO: Remove this command for OSSing. Google Drive extensions should be kept in internal version
.PHONY=build-web-frontend-viewer-mode
build-web-frontend-viewer-mode: ./web/**/*.ts ./web/**/*.html ./web/**/*.sass
	cd web &&NG_APP_VIEWER_MODE=true NG_APP_VERSION="$(VERSION)" NG_APP_GTAG_ID="$(GTAG_ID)" NG_APP_ENABLE_GOOGLE_DRIVE_DATA_LOADER="true" NG_APP_REPORT_BUG_URL=$(BUG_REPORT_URL) NG_APP_DOCUMENT_URL=$(DOCUMENT_URL) npx ng build --output-path ../dist

## Deploy subcommands

# TODO: Remove gcertstatus command for OSSing
.PHONY=deploy-precondition
deploy-precondition:
	./scripts/git-precondition.sh $(GIT_TAG_NAME)
	gcertstatus --check_remaining=15m

.PHONY=deploy-container-image
deploy-container-image: build-web ./pkg/**/*.go Dockerfile
	$(GCLOUD) builds submit --config=./cloudbuild.yaml --substitutions=_IMAGE_TAG="$(VERSION)"

.PHONY=deploy-container-image-beta
deploy-container-image-beta: build-web-beta ./pkg/**/*.go Dockerfile
	$(GCLOUD) builds submit --config=./cloudbuild-beta.yaml --substitutions=_IMAGE_TAG="$(VERSION)-beta"

.PHONY=deploy-analytics
deploy-analytics:
	docker build --file ./Dockerfile-analytics . --tag gcr.io/kubernetes-history-inspector/analytics:latest
	docker push gcr.io/kubernetes-history-inspector/analytics:latest
	$(GCLOUD) run deploy khi-analytics --image gcr.io/kubernetes-history-inspector/analytics:latest --region us-central1 --no-allow-unauthenticated

.PHONY=test-web
test-web:
	cd web && ng test --watch=false

.PHONY=test-go
test-go:
	go test ./...

.PHONY=clean-go-test-cache
clean-go-test-cache:
	go clean -testcache

.PHONY=lint-web
lint-web:
	cd web && npm run lint && npx stylelint "*/**.sass"

.PHONY=lint-go
lint-go:
	go vet ./...

.PHONY=format-go
format-go:
	gofmt -s -w .

.PHONY=format-web
format-web:
	cd web && npm run format

.PHONY=check-format-go
check-format-go:
	test -z `gofmt -l .`

.PHONY=check-format-web
check-format-web:
	cd web && npm run check-format

.PHONY=qa
qa: clean-go-test-cache test lint check-format-go check-format-web

.PHONY=coverage-web
coverage-web:
	cd web && ng test --code-coverage

.PHONY=coverage-go
coverage-go:
	go test -cover ./... -coverprofile=./go-cover.output
	go tool cover -html=./go-cover.output -o=go-cover.html
