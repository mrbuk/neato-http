FROM golang:1.16-alpine AS build

WORKDIR /go/src/neato-http
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -v ./...

# build a minimal image

FROM scratch

# the tls certificates:
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# the actual binary
COPY --from=build /go/bin/neato-http /go/bin/neato-http

EXPOSE 8080

ENTRYPOINT ["/go/bin/neato-http"]
