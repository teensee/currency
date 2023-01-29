FROM golang:1.18.0-alpine as build
RUN apk add alpine-sdk

RUN mkdir /currency
WORKDIR /currency
COPY ./ .

RUN go clean --modcache
RUN apk update && apk upgrade

RUN GOOS=linux GOARCH=amd64 go build -a -tags musl -a -installsuffix cgo -o app cmd/app/main.go

FROM golang:alpine
COPY --from=build /currency/app /

EXPOSE 8080
ENTRYPOINT ["/app"]