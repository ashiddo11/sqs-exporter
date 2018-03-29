package collector

import (
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "log"
        "github.com/aws/aws-sdk-go/service/sqs"
        "net/http"
        "strings"
        "fmt"
)

type MetricHandler struct{}

func (h MetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        queues := getQueues()
        for queue, attr := range queues {
                msgAvailable := *attr.Attributes["ApproximateNumberOfMessages"]
                msgDelayed := *attr.Attributes["ApproximateNumberOfMessagesDelayed"]
                msgNotvisible := *attr.Attributes["ApproximateNumberOfMessagesDelayed"]
                fmt.Fprintf(w, "sqs_messages_visible{\"queue\":%s} %+v\n", queue, msgAvailable)
                fmt.Fprintf(w, "sqs_messages_delayed{\"queue\":%s} %+v\n", queue, msgDelayed)
                fmt.Fprintf(w, "sqs_messages_not_visible{\"queue\":%s} %+v\n", queue, msgNotvisible)
        }
}

func getQueueName(url string) (queueName string) {
        queue := strings.Split(url, "/")
        queueName = queue[len(queue) -1 ]
        return
}

func getQueues() (queues map[string]*sqs.GetQueueAttributesOutput) {
        sess := session.Must(session.NewSession(&aws.Config{
                Region: aws.String("ap-southeast-1")}))
        client := sqs.New(sess)
        result, err := client.ListQueues(nil)
        if err != nil {
                log.Println("Error", err)
                return
        }

        queues = make(map[string]*sqs.GetQueueAttributesOutput)

        for _, urls := range result.QueueUrls {
                if urls == nil {
                        continue
                }
                params := &sqs.GetQueueAttributesInput{
                        QueueUrl: aws.String(*urls),
                        AttributeNames: []*string{
                                aws.String("ApproximateNumberOfMessages"),
                                aws.String("ApproximateNumberOfMessagesDelayed"),
                                aws.String("ApproximateNumberOfMessagesNotVisible"),
                        },
                }

                resp, _ := client.GetQueueAttributes(params)
                queueName := getQueueName(*urls)
                queues[queueName] = resp
        }
        return queues
}
