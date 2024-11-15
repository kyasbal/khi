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
build-web-internal: prepare-frontend ./web/**/*.ts ./web/**/*.html ./web/**/*.sass
	cd web && NG_APP_VERSION="$(VERSION)" npx ng build --output-path ../dist -c prod-internal

.PHONY=build-web-internal-beta
build-web-internal-beta: prepare-frontend ./web/**/*.ts ./web/**/*.html ./web/**/*.sass
	cd web && NG_APP_VERSION="beta-$(VERSION)@$(GIT_SHORT_HASH)" npx ng build --output-path ../dist --optimization=false -c dev-internal

.PHONY=build-web-internal-khi-ro
build-web-internal-khi-ro: prepare-frontend ./web/**/*.ts ./web/**/*.html ./web/**/*.sass
	cd web && NG_APP_VIEWER_MODE=true NG_APP_VERSION="$(VERSION)" npx ng build --output-path ../dist -c prod-internal


.PHONY=deploy-analytics
deploy-analytics:
	docker build --file ./Dockerfile-analytics . --tag gcr.io/kubernetes-history-inspector/analytics:latest
	docker push gcr.io/kubernetes-history-inspector/analytics:latest
	$(GCLOUD) run deploy khi-analytics --image gcr.io/kubernetes-history-inspector/analytics:latest --region us-central1 --no-allow-unauthenticated

.PHONY=remove-internal-code
remove-internal-code:
	./scripts/private/remove-internal-codes.sh

.PHONY=add-licenses
add-licenses:
	$(GOPATH)/bin/addlicense  -c "Google LLC" -l apache .