server:
	PORT=2565 AUTH_TOKEN="November 10, 2009" go run server.go

test:
	AUTH_TOKEN="November 10, 2009" go test -cover -v ./... --tags=unit

test-cover:
	AUTH_TOKEN="November 10, 2009" go test -coverprofile=coverage.out -v ./... --tags=unit
	go tool cover -html=coverage.out

docker-it-test-up:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

docker-it-test-down:
	docker-compose -f docker-compose.test.yml down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down