FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY devWhisper /app/devWhisper
COPY .env /app/.env

CMD ["/app/devWhisper"]
