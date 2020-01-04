FROM golang:latest AS build

WORKDIR /go/src
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
COPY . .
RUN make

FROM scratch

COPY --from=build /go/src/build/smartdns /smartdns

ENTRYPOINT ["/smartdns"]
