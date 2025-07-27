FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY devWhisper /usr/bin/devWhisper

CMD ["devWhisper"]