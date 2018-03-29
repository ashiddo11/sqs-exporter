FROM golang:1.10

WORKDIR /go/src/github.com/ashiddo11/sqs-exporter/

RUN go get github.com/aws/aws-sdk-go
 
COPY .  .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sqs-exporter .

FROM scratch

COPY --from=0 /go/src/github.com/ashiddo11/sqs-exporter/sqs-exporter /

EXPOSE 9434

CMD ["/sqs-exporter"]
