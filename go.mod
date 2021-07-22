module github.com/asyncapi/event-gateway

go 1.16

require (
	github.com/asyncapi/parser-go v0.3.1-0.20210701222435-43ab3e4b47d6
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/grepplabs/kafka-proxy v0.2.8
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.12.2
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pierrec/lz4 v2.6.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/xdg/scram v1.0.3 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210427231257-85d9c07bbe3a // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/grepplabs/kafka-proxy => github.com/smoya/kafka-proxy v0.2.9-0.20210623125118-a94cf71a065c
