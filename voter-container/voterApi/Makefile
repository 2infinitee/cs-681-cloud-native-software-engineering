SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build				Build the voter executable"
	@echo "	   run					Run the voter program from code"
	@echo "	   run-bin				Run the voter executable"
	@echo "	   load-db				Add sample data via curl"
	@echo "	   get-by-id			Get a voter by id pass id=<id> on command line"
	@echo "	   get-all				Get all voters"
	@echo "	   update-2				Update record 2, pass a new title in using title=<title> on command line"
	@echo "	   delete-all			Delete all voters"
	@echo "	   delete-by-id			Delete a voter by id pass id=<id> on command line"
	@echo "	   get-v2				Get all voters by done status pass done=<true|false> on command line"
	@echo "	   get-v2-all			Get all voters using version 2"
	@echo "	   build-amd64-linux	Build amd64/Linux executable"
	@echo "	   build-arm64-linux	Build arm64/Linux executable"

.PHONY: build
build:
	go build .

.PHONY: build-container
build-container:
	docker build --tag cs-t681-voter-api:v2 -f ./Dockerfile .

.PHONY: build-amd64-linux
build-amd64-linux:
	GOOS=linux GOARCH=amd64 go build -o ./voterApi-linux-amd64 .

.PHONY: build-arm64-linux
build-arm64-linux:
	GOOS=linux GOARCH=arm64 go build -o ./voterApi-linux-arm64 .

.PHONY: run
run:
	go run main.go

.PHONY: run-bin
run-bin:
	./voterApi

.PHONY: restore-db
restore-db:
	(cp ./data/voter-db.json.bak ./data/voter-db.json)

.PHONY: load-db
load-db:
	curl -d '{"voterId":1,"firstName":"yoyo","lastName":"ma","isDone":false,"voterHistory":[{"pollId":123,"voterId":1,"voteDate":"2006-01-02T15:04:05Z"},{"pollId":567,"voterId":1,"voteDate":"2006-01-02T15:04:05Z"}]}' -H "Content-Type: application/json" -X POST http://localhost:8080/voter
	curl -d '{"voterId":2,"firstName":"liam","lastName":"nelson","isDone":false,"voterHistory":[{"pollId":198,"voterId":2,"voteDate":"2006-01-02T15:04:05Z"},{"pollId":288,"voterId":2,"voteDate":"2006-01-02T15:04:05Z"}]}' -H "Content-Type: application/json" -X POST http://localhost:8080/voter
	curl -d '{"voterId":3,"firstName":"roger","lastName":"bowman","isDone":false,"voterHistory":[{"pollId":144,"voterId":3,"voteDate":"2006-01-02T15:04:05Z"},{"pollId":288,"voterId":3,"voteDate":"2006-01-02T15:04:05Z"}]}' -H "Content-Type: application/json" -X POST http://localhost:8080/voter
	curl -d '{"voterId":4,"firstName":"ken","lastName":"masters","isDone":false,"voterHistory":[{"pollId":554,"voterId":4,"voteDate":"2006-01-02T15:04:05Z"},{"pollId":285,"voterId":4,"voteDate":"2006-01-02T15:04:05Z"}]}' -H "Content-Type: application/json" -X POST http://localhost:8080/voter

.PHONY: update-2
update-2:
	curl -d '{ "voterId": 2, "firstName": "will", "done": false }' -H "Content-Type: application/json" -X PUT http://localhost:8080/voter

.PHONY: get-by-id
get-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/voter/$(voterId)

.PHONY: get-all
get-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/voter

.PHONY: delete-by-id
delete-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:8080/voter/2

.PHONY: delete-all
delete-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:8080/voter

.PHONY: get-v2
get-v2:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/v2/voter?done=$(done)

.PHONY: get-v2-all
get-v2-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:8080/v2/voter
