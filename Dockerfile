FROM golang:alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /src/app/cmd
RUN go build -o /main ./server/main.go
RUN go build -o /migrate ./migrate/main.go

FROM alpine:latest
RUN apk add --no-cache postgresql-client && rm -rf /var/cache/apk/*

COPY --from=builder /main /main
COPY --from=builder /migrate /migrate
COPY --from=builder /src/.env /.env

EXPOSE 8080

CMD ["/main"]
