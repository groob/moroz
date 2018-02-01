FROM alpine:3.6

RUN apk --update add \
    ca-certificates 

COPY ./build/linux/moroz /usr/bin/moroz

CMD ["moroz"]
