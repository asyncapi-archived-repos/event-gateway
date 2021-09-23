module github.com/asyncapi/event-gateway

go 1.16

require (
	github.com/Shopify/sarama v1.29.1
	github.com/asyncapi/parser-go v0.3.1-0.20210701222435-43ab3e4b47d6
	github.com/go-chi/chi/v5 v5.0.3
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grepplabs/kafka-proxy v0.2.8
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/olahol/melody v0.0.0-20180227134253-7bd65910e5ab
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/xdg/scram v1.0.3 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
)

replace (
	github.com/Shopify/sarama => github.com/smoya/sarama v1.29.1-0.20210922162934-6d18e39ca823
	github.com/grepplabs/kafka-proxy => github.com/smoya/kafka-proxy v0.2.9-0.20210623125118-a94cf71a065c
)
