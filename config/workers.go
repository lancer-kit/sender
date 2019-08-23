package config

const ServiceName = "sender"

const (
	WorkerAPIServer     = "api-server"
	WorkerAsyncAPIEmail = "async-api-email"
	WorkerAsyncAPISms   = "async-api-sms"
)

type Workers []string
