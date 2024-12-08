# DeBeAndo Agent for AWS/RDS

Database monitoring tool designed for small environments, adapted for Kubernetes and send metrics to InfluxDB.

## Image Description

This image is maintained by DeBeAndo and will be updated regularly on best-effort basis. The image is based on Alpine Linux and only contains the build result of this repository.

## Run

To run container:

```bash
docker run \
	--name debeando-agent-aws-rds \
	--env DEBUG=true \
	--env INTERVAL=10 \
	--env INFLUXDB_HOST="http://com-env-influxdb-observability-node01.aws.com" \
	--env INFLUXDB_TOKEN="abc123cde456==" \
	--env AWS_RDS_INSTANCE_IDENTIFIER="db1" \
	debeando/agent-aws-rds
```

## Environment Variables

When you start the `agent-aws-rds` image, you can adjust the configuration of the agent instance by passing one or more environment variables on the docker run command line.

- **DEBUG:** Enable debug mode with `true` value, by default value is `false`.
- **INTERVAL:** Interval time in second, by default value is `60`.
- **INFLUXDB_HOST:** The HTTP hostname or IP address.
- **INFLUXDB_TOKEN:** The authentication token for connecting to the InfluxDB instance.
- **AWS_RDS_INSTANCE_IDENTIFIER:** The DB instance name on AWS/RDS.
