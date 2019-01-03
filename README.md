# apcups-influx

[![Version Badges relate to](https://images.microbadger.com/badges/version/tnwhitwell/apcups-influx.svg)](https://microbadger.com/images/tnwhitwell/apcups-influx "Get your own version badge on microbadger.com") [![Image Size Badge](https://images.microbadger.com/badges/image/tnwhitwell/apcups-influx.svg)](https://microbadger.com/images/tnwhitwell/apcups-influx "Get your own image badge on microbadger.com") [![Built Commit Badge](https://images.microbadger.com/badges/commit/tnwhitwell/apcups-influx.svg)](https://microbadger.com/images/tnwhitwell/apcups-influx "Get your own commit badge on microbadger.com")

Available on dockerhub: [tnwhitwell/apcups-influx](https://hub.docker.com/r/tnwhitwell/apcups-influx)

## Description

This container will get metrics from apcupsd (via TCP) and publish them to an influxdb server at a specified interval

## Environmental Variables

This is a list of Environmental variables, and some sample values:

```sh
APCUPSD_URL=upshost:3551
INFLUXDB_DB=UPSData
INFLUXDB_URL=http://influxdb:8086
INFLUXDB_USER= # Leave blank / omit for no username
INFLUXDB_PASS= # Leave blank / omit for no password
READING_INTERVAL=60
```

## Samples to run

### `docker run`

`docker run -e APCUPSD_URL=upshost:3551 -e INFLUXDB_DB=UPSData -e INFLUXDB_URL=http://influxdb:8086 -e READING_INTERVAL=60 tnwhitwell/apcups-influx:latest`

### `docker-compose`

```yaml
version: '3'
services:
  influxdb:
    image: influxdb:latest
    restart: always
    volumes:
      - ${CONFIG}/influxdb:/var/lib/influxdb
  apcups-sender:
    image: tnwhitwell/apcups-influx:latest
    restart: always
    environment:
      APCUPSD_URL: upshost:3551
      READING_INTERVAL: 60
      INFLUXDB_DB: UPSData
      INFLUXDB_URL: http://influxdb:8086
```

`docker-compose up`
