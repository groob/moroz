FROM alpine:3.4

RUN apk --update add \
    ca-certificates 

COPY ./build/moroz-linux-amd64 /moroz

CMD ["/moroz"]

