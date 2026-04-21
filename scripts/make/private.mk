# private.mk
# This file contains make tasks used internally. This file will not be part of OSSed code.

GCLOUD_PROJECT="kubernetes-history-inspector"
GCLOUD=gcloud --project $(GCLOUD_PROJECT)

.PHONY=watch-web-internal
watch-web-internal: $(GENERATE_FRONTEND_DUMMY) ## Run frontend development server for Google internal production version
	cd web && npx ng serve -c dev-internal

.PHONY=watch-web-internal-khi-ro
watch-web-internal-khi-ro: $(GENERATE_FRONTEND_DUMMY) ## Run frontend development server for the read-only build of Google internal production version
	cd web && npx ng serve -c dev-internal-khi-ro


.PHONY=build-web-internal
build-web-internal: $(FRONTEND_ARTIFACT_FILES_DUMMY) ./web/**/*.ts ./web/**/*.html ./web/**/*.scss ## Build frontend for Google internal production version
	cd web && NG_APP_VERSION="$(VERSION)" npx ng build --output-path ../pkg/server/dist -c prod-internal

.PHONY=build-web-internal-beta
build-web-internal-beta: $(FRONTEND_ARTIFACT_FILES_DUMMY) ./web/**/*.ts ./web/**/*.html ./web/**/*.scss ## Build frontend for Google internal production version as a beta version
	cd web && NG_APP_VERSION="beta-$(VERSION)@$(GIT_SHORT_HASH)" npx ng build --output-path ../pkg/server/dist --optimization=false -c dev-internal

.PHONY=build-web-internal-khi-ro
build-web-internal-khi-ro: $(FRONTEND_ARTIFACT_FILES_DUMMY) ./web/**/*.ts ./web/**/*.html ./web/**/*.scss ## Build frontend for Google internal production version as a read-only build
	cd web && NG_APP_VIEWER_MODE=true NG_APP_VERSION="$(VERSION)" npx ng build --output-path ../pkg/server/dist -c prod-internal


.PHONY=deploy-analytics
deploy-analytics: ## Build the analytics proxy and deploy it on the Cloud Run
	podman build --platform linux/amd64 --file ./Dockerfile-analytics . --tag gcr.io/khi-internal/analytics:latest
	podman push gcr.io/khi-internal/analytics:latest
	$(GCLOUD) run deploy khi-analytics --image gcr.io/khi-internal/analytics:latest --region us-central1 --allow-unauthenticated

.PHONY=cleanup-khi-ro-versions
cleanup-khi-ro-versions: ## Clean-up the khi-ro test build version and keeps the latest 30 versions only
	gcloud app versions --project google.com:khi-ro list --format json | jq '[.[] |select(.id != "main")|select( .id != "beta")]|sort_by(.version.createTime)|reverse|.[].id' -r | tail -n +30 | xargs -I@ bash -c -x "gcloud app versions delete @ --project google.com:khi-ro"
