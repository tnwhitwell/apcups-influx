FROM golang:1.11.4-stretch as build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/app
ADD src/* /go/src/app/
RUN go get && go build -o /apc-influx

FROM scratch
ARG VCS_REF
ARG BUILD_DATE

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/tnwhitwell/apcups-influx" \
      org.label-schema.docker.cmd="docker run -e APCUPSD_URL=upshost:3551 -e INFLUXDB_DB=UPSData -e INFLUXDB_URL=http://influxdb:8086 -e READING_INTERVAL=60 tnwhitwell/apcups-influx" \
      org.label-schema.schema-version="1.0" \
      maintainer="tom@whi.tw"

COPY --from=build /apc-influx /apc-influx
CMD [ "/apc-influx" ]
