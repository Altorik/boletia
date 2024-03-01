FROM golang:alpine AS build

RUN apk add --update git
RUN apk --update add ca-certificates
WORKDIR /go/src/boletia
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/boletia cmd/api/main.go

FROM scratch

COPY --from=build /go/bin/boletia /go/bin/boletia
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/boletia"]
