export DATABASE_URL=postgres://qbxrkhsv:Ugq4c2uuZe6xAL6MZ2USrwc84Tzg4Uwq@john.db.elephantsql.com/qbxrkhsv
export DATABASE_DRIVER=postgres
export PORT=2565
export AUTH_TOKEN=November 10, 2009

server:
	go run server.go

test:
	go test -cover -v ./... --tags=unit

test-cover:
	go test -coverprofile=coverage.out -v ./... --tags=unit
	go tool cover -html=coverage.out

docker-it-test-up:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

docker-it-test-down:
	docker-compose -f docker-compose.test.yml down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down