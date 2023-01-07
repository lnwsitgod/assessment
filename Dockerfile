FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 AUTH_TOKEN="November 10, 2009" go test --tags=unit -v ./...

RUN go build -o ./out/assessment .

FROM alpine:3.16.2
COPY --from=build-base /app/out/assessment /app/assessment

CMD ["/app/assessment"]
