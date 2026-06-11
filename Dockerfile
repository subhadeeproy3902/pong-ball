FROM alpine:3.19
RUN apk add --no-cache ncurses
COPY paddle-ball /usr/local/bin/paddle-ball
ENTRYPOINT ["paddle-ball"]
CMD ["play"]