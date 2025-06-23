FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src .

RUN go build -o /discovery-app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /discovery-app /discovery-app

EXPOSE 5112

ENTRYPOINT ["/discovery-app"]