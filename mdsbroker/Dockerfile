FROM golang:1.22 as build


WORKDIR /go/src/app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /http_server http_server.go


FROM gcr.io/distroless/base
COPY --from=build /http_server /http_server
EXPOSE 8080

ENTRYPOINT ["/http_server"]
