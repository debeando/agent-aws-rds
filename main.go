package main

import (
	"context"
	"time"

	"github.com/debeando/go-common/env"
	"github.com/debeando/go-common/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var values = []string{
	"BinLogDiskUsage",
	"BurstBalance",
	"CPUUtilization",
	"CPUCreditUsage",
	"CPUCreditBalance",
	"CPUSurplusCreditBalance",
	"CPUSurplusCreditsCharged",
	"DatabaseConnections",
	"DiskQueueDepth",
	"FreeableMemory",
	"FreeStorageSpace",
	"LVMReadIOPS",
	"LVMWriteIOPS",
	"NetworkReceiveThroughput",
	"NetworkTransmitThroughput",
	"ReadIOPS",
	"ReadLatency",
	"ReadThroughput",
	"ReplicaLag",
	"SwapUsage",
	"WriteIOPS",
	"WriteLatency",
	"WriteThroughput",
	"NumVCPUs",
}

func main() {
	log.Info("Start DeBeAndo Agent for AWS RDS Metrics")

	if getDebug() {
		log.SetLevel(log.DebugLevel)
	}

	log.DebugWithFields("Environment Variables", log.Fields{
		"DEBUG":                       getDebug(),
		"INFLUXDB_BUCKET":             influxDB.Bucket,
		"INFLUXDB_HOST":               influxDB.Host,
		"INFLUXDB_PORT":               influxDB.Port,
		"INFLUXDB_TOKEN":              influxDB.Token,
		"INTERVAL":                    getInterval(),
		"AWS_RDS_INSTANCE_IDENTIFIER": AWSRDSInstanceIdentifier,
	})

	influxDB.New()
	defer influxDB.Close()

	for {
		metrics := Metrics{}
		queries := []types.MetricDataQuery{}

		awsConfig, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Error(err.Error())
		}

		cwClient := cloudwatch.NewFromConfig(awsConfig)

		for _, value := range values {
			queries = append(queries,
				types.MetricDataQuery{
					Id: aws.String("id" + value),
					MetricStat: &types.MetricStat{
						Metric: &types.Metric{
							Namespace:  aws.String("AWS/RDS"),
							MetricName: aws.String(value),
							Dimensions: []types.Dimension{
								{
									Name:  aws.String("DBInstanceIdentifier"),
									Value: aws.String(AWSRDSInstanceIdentifier),
								},
							},
						},
						Period: aws.Int32(60),
						Stat:   aws.String("Average"),
					},
				})
		}

		input := &cloudwatch.GetMetricDataInput{
			EndTime:           aws.Time(time.Unix(time.Now().Unix(), 0)),
			StartTime:         aws.Time(time.Unix(time.Now().Add(time.Duration(-2)*time.Minute).Unix(), 0)),
			MetricDataQueries: queries,
		}

		result, err := cwClient.GetMetricData(context.TODO(), input)
		if err != nil {
			log.Error(err.Error())
		}

		metric := Metric{}
		metric.Measurement = "aws_rds"
		metric.AddTag(Tag{
			Name:  "server",
			Value: AWSRDSInstanceIdentifier,
		})

		for _, r := range result.MetricDataResults {
			if len(r.Values) > 0 {
				metric.AddField(Field{
					Name:  *r.Label,
					Value: r.Values[0],
				})
			}
		}

		metrics.Add(metric)

		if metrics.Count() > 0 {
			influxDB.Write(metrics)
		}

		metrics.Reset()
		log.Debug("Wait until next collect metrics.")
		time.Sleep(getInterval())
	}
}

func getDebug() bool {
	return env.GetBool("DEBUG", true)
}

func getInterval() time.Duration {
	return time.Duration(env.GetInt("INTERVAL", 60)) * time.Second
}
