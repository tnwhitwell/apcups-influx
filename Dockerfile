FROM golang:alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

RUN adduser -D -g '' appuser

ADD src/* ${GOPATH}/src/app/
WORKDIR ${GOPATH}/src/app

RUN go get -d -v
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/apc-influx

FROM scratch
ARG VCS_REF
ARG BUILD_DATE

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/tnwhitwell/apcups-influx" \
      org.label-schema.docker.cmd="docker run -e APCUPSD_URL=upshost:3551 -e INFLUXDB_DB=UPSData -e INFLUXDB_URL=http://influxdb:8086 -e READING_INTERVAL=60 tnwhitwell/apcups-influx" \
      org.label-schema.docker.params="APCUPSD_URL=url of apcupsd instance,INFLUXDB_DB=name of infuxdb database to use,INFLUXDB_USER=username for influxdb,INFLUXDB_PASS=password for influxdb,INFLUXDB_URL=URL of the influxdb instance,READING_INTERVAL=integer number of seconds between readings are taken," \
      org.label-schema.schema-version="1.0" \
      maintainer="tom@whi.tw"

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /go/bin/apc-influx /go/bin/apc-influx

USER appuser

ENTRYPOINT [ "/go/bin/apc-influx" ]
