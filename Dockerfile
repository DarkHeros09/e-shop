# Build stage
FROM golang:1.18.1-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN curl -L https://github.com/eficode/wait-for/releases/download/v2.2.3/wait-for -o wait-for.sh
RUN chmod +x wait-for.sh

# Run stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/wait-for.sh .
COPY app.env .
COPY start.sh .
COPY db/migration ./migration
EXPOSE 8080
ENTRYPOINT [ "/app/start.sh" ]
CMD [ "/app/main" ]