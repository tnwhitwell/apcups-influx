FROM golang:1.11.4-stretch as build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/app
ADD src/* /go/src/app/
RUN go get && go build -o /apc-influx

FROM scratch

COPY --from=build /apc-influx /apc-influx
CMD [ "/apc-influx" ]
