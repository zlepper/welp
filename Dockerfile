FROM golang:1.10 as builder

WORKDIR /go/src/github.com/zlepper/welp

COPY . .

RUN go get -d -v .
RUN CGO_ENABLED=0 go build -o welp

FROM alpine
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/zlepper/welp/welp /usr/bin/welp

EXPOSE 8080

VOLUME /db /storage

CMD ["/usr/bin/welp"]