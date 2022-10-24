package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gomarkdown/markdown"
)

func uploadToS3(svc *s3.S3, region string, bucket string, key string, body []byte) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Println("Error uploading to S3!\n\t" + err.Error())
		return (err)
	}
	fmt.Println(resp)
	return nil
}

func main() {
	// sort the command line arguments
	var bucket, key, region, filename string
	flag.StringVar(&bucket, "bucket", "my-bucket", "S3 bucket name")
	flag.StringVar(&key, "key", "key.html", "S3 key name")
	flag.StringVar(&region, "region", "us-west-1", "AWS region")
	flag.StringVar(&filename, "file", "README.md", "Markdown file to convert to HTML and Publish")
	flag.Parse()

	fmt.Println("Formatting Markdown to HTML")
	dat, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file!\n\t" + err.Error())
		return
	}
	md := []byte(string(dat))
	body := markdown.ToHTML(md, nil, nil)
	fmt.Println(string(body))

	// Upload the HTML to S3
	fmt.Println("Uploading to S3...")

	// Create a new session for S3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		fmt.Println("Error creating session!\n\t" + err.Error())
		return
	}

	fmt.Println()

	svc := s3.New(sess)

	err = uploadToS3(svc, region, bucket, key, body)

	if err != nil {
		return
	}

	fmt.Println("Done! See the page at\nhttp://" + bucket + ".s3-website-" + region + ".amazonaws.com/" + key)
}
