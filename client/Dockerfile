FROM golang:1.22 as build


WORKDIR /go/src/app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /http_client http_client.go


FROM gcr.io/distroless/base
COPY --from=build /http_client /http_client
EXPOSE 8080

ENTRYPOINT ["/http_client"]
