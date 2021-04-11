FROM golang:1.16-alpine as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o crawler cmd/crawler/main.go
RUN go build -o notifier cmd/notifier/main.go
RUN go build -o scheduler cmd/scheduler/main.go

# PRODUCTION
FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /build/crawler /app/
COPY --from=builder /build/notifier /app/
COPY --from=builder /build/scheduler /app/

WORKDIR /app

EXPOSE 2112

ENTRYPOINT ["./crawler"]
CMD [""]