REGISTRY=tools-docker-local.artifactory.swisscom.com
IMAGE=swisscom/bitbucket-approver-bot
VERSION=$(shell ./version.sh)
TAG=$(VERSION)

.PHONY: clean, docker-build, docker-run, docker-push

build:
	CGO_ENABLED=0 go build \
		-o ./approve-bot \
		-ldflags="-X 'main.version=$(VERSION)'"
clean:
	rm ./approve-bot

docker-build:
	docker build \
		--build-arg "VERSION=$(VERSION)" \
		-t "$(REGISTRY)/$(IMAGE):$(TAG)" \
		.
	docker tag "$(REGISTRY)/$(IMAGE):$(TAG)" "$(REGISTRY)/$(IMAGE):latest"

docker-run:
	docker run --rm \
		--name bitbucket-approver-bot \
		-e "BITBUCKET_USERNAME=$$BITBUCKET_USERNAME" \
		-e "BITBUCKET_PASSWORD=$$BITBUCKET_PASSWORD" \
		-e "BITBUCKET_ENDPOINT=$$BITBUCKET_ENDPOINT" \
		-e "BITBUCKET_AUTHOR_FILTER=$$BITBUCKET_AUTHOR_FILTER" \
		-e "BITBUCKET_ADD_COMMENT=$$BITBUCKET_ADD_COMMENT" \
		"$(REGISTRY)/$(IMAGE):$(TAG)"


docker-push:
	docker push "$(REGISTRY)/$(IMAGE):$(TAG)"
	docker push "$(REGISTRY)/$(IMAGE):latest"