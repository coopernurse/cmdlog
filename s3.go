package cmdlog

import (
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"time"
)

const (
	dateTimeFmt string = "2006-01-02T15:04:05"
	dateFmt     string = "2006-01-02"
)

func EnsureBucket(accessKey, secretKey, regionName, bucketName string, acl s3.ACL) (*s3.Bucket, error) {
	var region aws.Region
	var ok bool
	if regionName == "" {
		region = aws.USEast
	} else {
		region, ok = aws.Regions[regionName]
		if !ok {
			return nil, fmt.Errorf("cmdlog: Unknown AWS region: %s", regionName)
		}
	}

	auth, err := aws.GetAuth(accessKey, secretKey, "", time.Time{})
	if err != nil {
		return nil, err
	}

	s := s3.New(auth, region)
	bucket := s.Bucket(bucketName)
	err = bucket.PutBucket(acl)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

type S3Logger struct {
	Bucket      *s3.Bucket
	Perm        s3.ACL
	Options     s3.Options
	ContentType string
}

func (s S3Logger) Log(result *Result) error {
	return s.Bucket.Put(path(result), []byte(data(result)), s.ContentType, s.Perm, s.Options)
}

func (s S3Logger) URL(result *Result) string {
	return s.Bucket.URL(path(result))
}

func path(result *Result) string {
	return fmt.Sprintf("%s/%s-%s.txt", result.StartDate(dateFmt), result.Name, result.StartDate(dateTimeFmt))
}

const dataTmpl string = `            Job: %s
           Date: %s
Elapsed seconds: %.2f
           Host: %s
       Exit Str: %s

STDOUT:
%s

STDERR:
%s`

func data(result *Result) string {
	return fmt.Sprintf(dataTmpl, result.Name, result.StartDate(dateTimeFmt), result.ElapsedSeconds(),
		result.Host, result.ExitStr, string(result.Stdout), string(result.Stderr))
}
