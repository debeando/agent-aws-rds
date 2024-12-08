package main

import (
	"github.com/debeando/go-common/env"
)

var AWSRDSInstanceIdentifier string

func init() {
	AWSRDSInstanceIdentifier = env.Get("AWS_RDS_INSTANCE_IDENTIFIER", "")
}
