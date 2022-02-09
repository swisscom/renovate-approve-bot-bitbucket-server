REGISTRY=tools-docker-local.artifactory.swisscom.com
IMAGE=swisscom/bitbucket-approver-bot
TAG=$(shell head -n1 VERSION)

.PHONY: docker-build, docker-run, docker-push

docker-build:
	docker build . \
		-t "$(REGISTRY)/$(IMAGE):$(TAG)"
	docker tag "$(REGISTRY)/$(IMAGE):$(TAG)" "$(REGISTRY)/$(IMAGE):latest"

docker-run:
	docker run --rm \
		--name bitbucket-approver-bot \
		-e "BITBUCKET_USERNAME=$$BITBUCKET_USERNAME" \
		-e "BITBUCKET_PASSWORD=$$BITBUCKET_PASSWORD" \
		-e "BITBUCKET_ENDPOINT=$$BITBUCKET_ENDPOINT" \
		-e "BITBUCKET_AUTHOR_FILTER=$$BITBUCKET_AUTHOR_FILTER" \
		"$(REGISTRY)/$(IMAGE):$(TAG)"


docker-push:
	docker push "$(REGISTRY)/$(IMAGE):$(TAG)"
	docker push "$(REGISTRY)/$(IMAGE):latest"