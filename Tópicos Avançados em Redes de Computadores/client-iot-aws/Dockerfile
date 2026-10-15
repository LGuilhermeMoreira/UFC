FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o server ./cmd/api
RUN chmod +x /app/server


FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/server /app/
COPY --from=builder /app/web /web
COPY --from=builder /app/certs /certs
EXPOSE 8080
CMD [ "/app/server" ]
