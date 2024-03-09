SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "	Targets:"
	@echo "		setup - This builds a working setup for docker compose in this lab. It uses a container on my personal docker repo (arm64 based)."
	@echo "		clean - This will clean out all images by running a docker stop and remove all container. **BE CAUTIOUS**"
	@echo "		build container - This will build the container based on your local machine's architecture."
	@echo "		init-redis - This will build a setup from a local container you have to built and deploy it to docker compose. **Build the container first!!!***"


.PHONY: build-container
build-container:
	docker build --tag cs-t681-voter-api:v3 -f ./voterApi/Dockerfile voterApi

.PHONY: runVoterApi
runVoterApi:
	docker run -it -rm -name jy-voter-api -p 8080:8080 cs-t681-voter-api:v1

.PHONY: basic
basic: clean setup-one load-db

.PHONY: network
network: clean setup-two load-db get-redis-cli

.PHONE: volume
volume: clean setup-three load-db

.PHONE: init-redis
init-redis: clean setup-four

.PHONE: setup
setup: build-container setup-five get-redis-cli

.PHONY: setup-one
setup-one:
	docker compose -f ./redis/1-basic/redis-docker-compose.yaml up -d

.PHONY: setup-two
setup-two:
	docker compose -f ./redis/2-network/redis-docker-compose.yaml up -d

.PHONY: setup-three
setup-three:
	docker compose -f ./redis/3-volume/redis-docker-compose.yaml up -d

.PHONY: setup-four
setup-four:
	docker compose -f ./redis/4-init-redis/redis-docker-compose.yaml up -d

.PHONY: setup-fice
setup-five:
	docker compose -f ./redis/4-init-redis/redis-docker-compose-pull-only.yaml up -d

.PHONY: get-redis-cli
get-redis-cli:
	@echo -e ""
	@echo -e ""
	@echo -e "Run below command to connect to redis backend"
	@echo -e "	./redis/2-network/redis-cli.sh"

.PHONY: load-db
load-db:
	./redis/1-basic/loadcache.sh

.PHONY: clean
clean:
	-docker stop $$(docker ps -aq)
	-docker rm $$(docker ps -aq)