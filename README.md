# Simple AWS SQS exporter
A prometheus SQS metrics exporter

### Metrics

| Metric  | Labels | Description |
| ------  | ------ | ----------- |
| sqs_messages_visible | Queue Name | Number of messages available |
| sqs_messages_delayed | Queue Name | Number of messages delayed |
| sqs_messages_not_visible | Queue Name | Number of messages in flight |

For more information check [AWS SQS Documentation](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-message-attributes.html)

## Configuration
Credentials to AWS are provided in the following order:

- Environment variables (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY)

- Shared credentials file (~/.aws/credentials)

- IAM role for Amazon EC2

For more information check [AWS SDK Documentation] (https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

### AWS IAM permissions
The app needs sqs list and read access to the sqs policies
```{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "sqs:ListQueues",
                "sqs:GetQueueUrl",
                "sqs:ListDeadLetterSourceQueues",
                "sqs:ReceiveMessage",
                "sqs:GetQueueAttributes",
                "sqs:ListQueueTags"
            ],
            "Resource": "*"
        }
    ]
}
```

## Running
**You need to specify the region you to connect to**
Running on an ec2 machine using IAM roles:
`docker run -e AWS_REGION=<region> -d -p 9434:9434 ashiddo11/sqs-exporter`

Or running it externally:
```docker run -d -p 9384:9384 -e AWS_ACCESS_KEY_ID=<access_key> -e AWS_SECRET_ACCESS_KEY=<secret_key> -e AWS_REGION=<region>  ashiddo11/sqs-exporter```
