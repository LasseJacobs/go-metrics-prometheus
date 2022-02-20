FROM golang:1.17-alpine as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /usr/local/go/src/github.com/LasseJacobs/go-metrics-prometheus
# Pulling dependencies
COPY ./go.* ./
RUN go mod download

# Building stuff
COPY ./prometheus ./prometheus
COPY _test/main.tgo ./main.go
RUN go build -o prom-test .

FROM scratch
# need shell for this; something to investigate
#RUN adduser -D -u 1000 runman

COPY --from=build /usr/local/go/src/github.com/LasseJacobs/go-metrics-prometheus/prom-test prom-test

EXPOSE 8080

#USER runman
ENTRYPOINT ["./prom-test"]