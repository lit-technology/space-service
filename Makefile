PORT?=8000
PACKAGE:=github.com/philip-bui/space-service
APPLICATION:=space
COVERAGE:=coverage.out
POSTGRES_DOCKER:=postgres:11.1
POSTGRES_NAME:=postgres-space
POSTGRES_USER:=philip
POSTGRES_PW:=KFC
POSTGRES_DB:=space

help: ## Display this help screen
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

godoc: ## Open HTML API Documentation
	echo "localhost:${PORT}/pkg/${PACKAGE}"
	godoc -http=:${PORT}

zip: ## Zip Package for ELB
	env GOOS=linux GOARCH=amd64 go build
	mv space-service application
	zip application.zip application

deploy: zip ## Deploy to ELB
	eb deploy

setup: mod protos postgres tables

mod: # Downloads dependencies
	export GOMODULES111=on
	go get ./...

test: ## Run Tests
	go test -coverprofile=coverage.out ./...

benchmark: ## Run Benchmark Tests
	go test -v -bench=. ./..

coverage: test ## Open HTML Test Coverage Report
	go tool cover -html=${COVERAGE}

proto: ## Generate Protobuf files
	go install github.com/gogo/protobuf/protoc-gen-gofast
	protoc -I protos/ protos/*.proto --gofast_out=grpc:protos

json: ## Generate JSON files
	go install github.com/francoispqt/gojay/gojay
	gojay -s protos/ -p True -t Post -pkg ${APPLICATION} -o protos/${APPLICATION}.json.go

postgres-stop: ## Stop Postgres Docker
	docker ps -aq --filter name=${POSTGRES_NAME} | xargs docker stop

postgres: ## Run Postgres Docker
	docker run --name ${POSTGRES_NAME} -e POSTGRES_DB=${POSTGRES_DB} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PW=${POSTGRES_PW} -d -p 5432:5432 ${POSTGRES_DOCKER}

tables: ## Insert Postgres Tables
	export PGPASSWORD=${POSTGRES_PW}
	docker cp ./resources/ ${POSTGRES_NAME}:/
	for sql in `ls -A1 ./resources`; do \
		docker exec -it ${POSTGRES_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -f resources/$$sql; \
	done;
