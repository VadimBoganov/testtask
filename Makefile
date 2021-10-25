build:
	go mod download && CGO_ENABLED=0 GOOS=windows go build -o ./.bin/app ./cmd/main.go

run:
	docker-compose up --remove-orphans

export HTTP_SERVER_CONTAINER_NAME=apache_test
export MONGODB_CONTAINER_NAME=mongodb_test

test:
	make test.unit
	make test.integration
	make test.coverage

test.unit:
	go test -v -short -coverprofile=cover.out ./...

test.integration:
	docker run --rm -d -p 8080:80 --name $$HTTP_SERVER_CONTAINER_NAME httpd
	docker cp test_file.csv $$HTTP_SERVER_CONTAINER_NAME:/usr/local/apache2/htdocs/
	docker run --rm -d -p 27018:27017 --name $$MONGODB_CONTAINER_NAME mongo

	GIN_MODE=release go test -v -run Integration -coverprofile=cover.out ./...

	docker stop $$MONGODB_CONTAINER_NAME
	docker stop $$HTTP_SERVER_CONTAINER_NAME

test.coverage:
	go tool cover -func=cover.out