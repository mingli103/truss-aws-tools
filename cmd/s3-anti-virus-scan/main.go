package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/trussworks/truss-aws-tools/internal/aws/session"
)

func main() {
	var bucket, profile, region string
	flag.StringVar(&bucket, "bucket", "", "The S3 bucket to get the size.")
	flag.StringVar(&region, "region", "", "The AWS region to use.")
	flag.StringVar(&profile, "profile", "", "The AWS profile to use.")
	flag.Parse()
	if bucket == "" {
		flag.PrintDefaults()
		return
	}
	s3Client := makeS3Client(region, profile)
	bucketregion, err := getBucketRegion(s3Client, bucket)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bucketregion)
}

// makeS3Client makes an S3 client
func makeS3Client(region, profile string) *s3.S3 {
	sess := session.MustMakeSession(region, "")
	c := s3.New(sess)
	return c
}

func getBucketRegion(c *s3.S3, bucket string) (string, error) {
	res, err := c.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", err
	}
	// Ugh. The S3 API returns inconsistent responses; namely, us-east-1
	// and eu-west-1 return non-standard values that need to be converted
	// into the standard region codes.
	// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTBucketGETlocation.html
	var region string
	if res.LocationConstraint != nil {
		region = *res.LocationConstraint
	} else {
		region = ""
	}
	switch region {
	case "":
		return "us-east-1", nil
	case "EU":
		return "eu-west-1", nil
	default:
		return region, nil
	}
}
