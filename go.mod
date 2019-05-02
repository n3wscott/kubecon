module github.com/n3wscott/kubecon

go 1.12

require (
	github.com/bradfitz/gomemcache v0.0.0-20190329173943-551aad21a668
	github.com/cloudevents/sdk-go v0.6.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.0
	github.com/kelseyhightower/envconfig v1.3.0
	pack.ag/amqp v0.11.0
)

replace github.com/cloudevents/sdk-go => github.com/n3wscott/sdk-go v0.0.0-20190501173006-4c0686d5867e
