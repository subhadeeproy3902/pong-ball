FROM alpine:3.19
RUN apk add --no-cache ncurses
COPY pong-ball /usr/local/bin/pong-ball
ENTRYPOINT ["pong-ball"]
CMD ["play"]