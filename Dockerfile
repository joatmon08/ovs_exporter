FROM golang:1.7

COPY . /ovs_exporter

RUN cd /ovs_exporter && go build ovs_exporter.go

CMD ovs_exporter