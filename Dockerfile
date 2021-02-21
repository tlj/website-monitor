FROM golang:1.16-alpine as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o website-monitor .

# PRODUCTION
FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /build/website-monitor /app/

WORKDIR /app

EXPOSE 2112

ENTRYPOINT ["./website-monitor"]
CMD [""]