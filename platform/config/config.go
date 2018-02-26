package config

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

func AWSConfig() *aws.Config {
	region := aws.String(Env.Aws.Region)
	if Env.Name == LocalEnv {
		return &aws.Config{
			Region: aws.String(Env.Aws.Region),
		}
	}

	ec2m := ec2metadata.New(session.New(), &aws.Config{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	})

	cf := &aws.Config{
		Region: region,
	}

	cr := credentials.NewCredentials(&ec2rolecreds.EC2RoleProvider{
		Client: ec2m,
	})
	cf.Credentials = cr

	return cf
}
