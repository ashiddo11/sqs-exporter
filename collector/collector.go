package collector

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"net/http"
	"strings"
)

type MetricHandler struct{}

func (h MetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	queues := getQueues()
	listQueueTags := getTags()
	for queue, attr := range queues {
		msgAvailable := *attr.Attributes["ApproximateNumberOfMessages"]
		msgDelayed := *attr.Attributes["ApproximateNumberOfMessagesDelayed"]
		msgNotvisible := *attr.Attributes["ApproximateNumberOfMessagesDelayed"]
		tags := ""
		for key, value := range listQueueTags[queue].Tags {
			tags += "," + key + "=\"" + *value + "\""
		}
		fmt.Fprintf(w, "sqs_messages_visible{queue_name=\"%s\"%s} %+v\n", queue, tags, msgAvailable)
		fmt.Fprintf(w, "sqs_messages_delayed{queue_name=\"%s\"%s} %+v\n", queue, tags, msgDelayed)
		fmt.Fprintf(w, "sqs_messages_not_visible{queue_name=\"%s\"%s} %+v\n", queue, tags, msgNotvisible)
	}
}

func getQueueName(url string) (queueName string) {
	queue := strings.Split(url, "/")
	queueName = queue[len(queue)-1]
	return
}

func getTags() (tags map[string]*sqs.ListQueueTagsOutput) {
	sess := session.Must(session.NewSession())
	client := sqs.New(sess)
	result, err := client.ListQueues(nil)
	if err != nil {
		log.Fatal("Error ", err)
	}

	tags = make(map[string]*sqs.ListQueueTagsOutput)

	if result.QueueUrls == nil {
		log.Println("Couldnt find any queues in region:", *sess.Config.Region)
	}
	for _, urls := range result.QueueUrls {
		params := &sqs.ListQueueTagsInput{
			QueueUrl: aws.String(*urls),
		}

		resp, _ := client.ListQueueTags(params)
		queueName := getQueueName(*urls)
		tags[queueName] = resp
	}
	return tags
}

func getQueues() (queues map[string]*sqs.GetQueueAttributesOutput) {
	sess := session.Must(session.NewSession())
	client := sqs.New(sess)
	result, err := client.ListQueues(nil)
	if err != nil {
		log.Fatal("Error ", err)
	}

	queues = make(map[string]*sqs.GetQueueAttributesOutput)

	if result.QueueUrls == nil {
		log.Println("Couldnt find any queues in region:", *sess.Config.Region)
	}
	for _, urls := range result.QueueUrls {
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
