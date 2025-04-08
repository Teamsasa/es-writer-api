FROM golang:alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR ./app/cmd
RUN go build -o /main ./server/main.go
RUN go build -o /migrate ./migrate/main.go

FROM alpine:latest
RUN apk add --no-cache postgresql-client && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /main ./main
COPY --from=builder /migrate ./migrate
COPY --from=builder /src/.env ./.env
COPY --from=builder /src/app/internal/usecase/prompts/es_generation.txt ./prompts/es_generation.txt
COPY --from=builder /src/app/internal/usecase/prompts/extract_questions.txt ./prompts/extract_questions.txt

EXPOSE 8080

CMD ["./main"]
