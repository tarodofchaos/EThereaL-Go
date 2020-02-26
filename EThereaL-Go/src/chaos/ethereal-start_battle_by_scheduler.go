package main

import (
	"fmt"
	"os"
	"strings"
	"context"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-lambda-go/lambda"
)


func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}

type S3Event struct {
        Message string `json:"scheduledBattle"`
}

func HandleRequest(ctx context.Context, result S3Event) (string, error) {
        sess, error := session.NewSession(&aws.Config{
	    Region: aws.String("eu-west-1")},
		)
		if error != nil {
		    exitErrorf("Error when getting session: " + error.Error())
		}
		
		bucketName := "ethereal-app"
		schedulePrefix := "private/battles"
		triggerPrefix := "private/army"
		
		s3Svc := s3.New(sess)
		resp, error := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket: aws.String(bucketName),
				Prefix: aws.String(schedulePrefix),
				})
		if error != nil {
		    exitErrorf("Error when listing objects: " + error.Error())
		}
		
		for _, item := range resp.Contents {
		    fmt.Println("Name:         ", *item.Key)
		    fmt.Println("Size:         ", *item.Size)
		    fmt.Println("")
		    if *item.Size > 0 {
			    resp, error := s3Svc.CopyObject(&s3.CopyObjectInput{
					Bucket: aws.String(bucketName),
					CopySource: aws.String(bucketName + "/" + *item.Key),
					Key: aws.String(strings.Replace(*item.Key, schedulePrefix, triggerPrefix, 1)),
				})	
			    if resp != nil {
			    	
			    }
			    if error != nil {
				    exitErrorf("Error when copying objects: " + error.Error())
				}
			    s3Svc.DeleteObject(&s3.DeleteObjectInput{
				    Bucket: aws.String(bucketName),
			    	Key: aws.String(*item.Key),	
		    	})
		    }
		    
		}
        
        return fmt.Sprintf("%s armies have been sent to battle!!", *item.Size-1 ), nil
}

func main() {
	lambda.Start(HandleRequest)
}