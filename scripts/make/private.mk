# private.mk
# This file contains make tasks used internally. This file will not be part of OSSed code.

GCLOUD_PROJECT="kubernetes-history-inspector"
GCLOUD=gcloud --project $(GCLOUD_PROJECT)


.PHONY=watch-web-internal
watch-web-internal: prepare-frontend
	cd web && NG_APP_BACKEND_URL_PREFIX="http://localhost:8080" ng serve -c dev-internal

.PHONY=watch-web-internal-khi-ro
watch-web-internal-khi-ro: prepare-frontend
	cd web && NG_APP_BACKEND_URL_PREFIX="http://localhost:8080" ng serve -c dev-internal-khi-ro


.PHONY=build-web-internal
build-web-internal: prepare-frontend ./web/**/*.ts ./web/**/*.html ./web/**/*.scss
	cd web && NG_APP_VERSION="$(VERSION)" npx ng build --output-path ../dist -c prod-internal

.PHONY=build-web-internal-beta
build-web-internal-beta: prepare-frontend ./web/**/*.ts ./web/**/*.html ./web/**/*.scss
	cd web && NG_APP_VERSION="beta-$(VERSION)@$(GIT_SHORT_HASH)" npx ng build --output-path ../dist --optimization=false -c dev-internal

.PHONY=build-web-internal-khi-ro
build-web-internal-khi-ro: prepare-frontend ./web/**/*.ts ./web/**/*.html ./web/**/*.scss
	cd web && NG_APP_VIEWER_MODE=true NG_APP_VERSION="$(VERSION)" npx ng build --output-path ../dist -c prod-internal


.PHONY=deploy-analytics
deploy-analytics:
	podman build --platform linux/amd64 --file ./Dockerfile-analytics . --tag gcr.io/khi-internal/analytics:latest
	podman push gcr.io/khi-internal/analytics:latest
	$(GCLOUD) run deploy khi-analytics --image gcr.io/khi-internal/analytics:latest --region us-central1 --allow-unauthenticated

.PHONY=generate-oss-repository
generate-oss-repository:
	./scripts/private/generate-oss-repository.sh

.PHONY=cleanup-khi-ro-versions
cleanup-khi-ro-versions:
	gcloud app versions --project google.com:khi-ro list --format json | jq '[.[] |select(.id != "main")|select( .id != "beta")]|sort_by(.version.createTime)|reverse|.[].id' -r | tail -n +30 | xargs -I@ bash -c -x "gcloud app versions delete @ --project google.com:khi-ro"