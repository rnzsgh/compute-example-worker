package cloud

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsc *sqs.SQS

func init() {
	sqsc = sqs.New(session.New())
}

// SqsDeleteMessage removes the message from the queue.
func SqsDeleteMessage(queue, receipt string) error {

	if len(queue) == 0 {
		panic(fmt.Sprintf("Empty queue to SqsSendMessage for queue: %s", queue))
	}

	if len(receipt) == 0 {
		panic(fmt.Sprintf("Empty receipt to SqsSendMessage for queue: %s", queue))
	}

	if _, err := sqsc.DeleteMessage(
		&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queue),
			ReceiptHandle: aws.String(receipt),
		},
	); err != nil {
		return fmt.Errorf("Failed to receive message from queue: %s - reason: %v", queue, err)
	}

	return nil
}

// SqsReceiveOneMessage pulls a single available message off the queue or returns
// an empty string. If the timeout or wait are negative this function panics. Set
// timeout or wait to zero to avoid setting.
func SqsReceiveOneMessage(queue string, timeout, wait int64) (string, string, error) {

	if len(queue) == 0 {
		panic(fmt.Sprintf("Empty queue to SqsSendMessage for queue: %s", queue))
	}

	if timeout < 0 {
		panic(fmt.Sprintf("Timeout is less than zero for SqsReceiveOneMessage for queue: %s", queue))
	}

	if wait < 0 {
		panic(fmt.Sprintf("Wait is less than zero for SqsReceiveOneMessage for queue: %s", queue))
	}

	var t *int64
	if timeout > 0 {
		t = aws.Int64(timeout)
	}

	var w *int64
	if wait > 0 {
		w = aws.Int64(wait)
	}

	res, err := sqsc.ReceiveMessage(
		&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queue),
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   t,
			WaitTimeSeconds:     w,
		},
	)

	if err != nil {
		return "", "", fmt.Errorf("Failed to receive message from queue: %s - reason: %v", queue, err)
	}

	if len(res.Messages) == 0 {
		return "", "", nil
	}

	return aws.StringValue(res.Messages[0].Body), aws.StringValue(res.Messages[0].ReceiptHandle), nil
}

// SqsSendMessage sends the message to the queue and returns the
// message id and an error if anything goes wrong. If the len of the payload
// or queue is zero, this function panics.
func SqsSendMessage(ctx aws.Context, queue, payload string) (string, error) {

	if len(queue) == 0 {
		panic(fmt.Sprintf("Empty queue to SqsSendMessage for queue: %s", queue))
	}

	if len(payload) == 0 {
		panic(fmt.Sprintf("Empty msg payload to SqsSendMessage for queue: %s", queue))
	}

	res, err := sqsc.SendMessageWithContext(
		ctx,
		&sqs.SendMessageInput{
			MessageBody: aws.String(payload),
			QueueUrl:    aws.String(queue),
		},
	)

	if err != nil {
		return "", fmt.Errorf("Failed to send message to queue: %s - reason: %v", queue, err)
	}

	return aws.StringValue(res.MessageId), nil
}
