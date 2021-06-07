module github.com/ONSdigital/dp-collection-api

go 1.16

replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible

require (
	github.com/ONSdigital/dp-component-test v0.3.0
	github.com/ONSdigital/dp-healthcheck v1.0.5
	github.com/ONSdigital/dp-mongodb v1.5.0
	github.com/ONSdigital/dp-net v1.0.12
	github.com/ONSdigital/log.go/v2 v2.0.0
	github.com/benweissmann/memongo v0.1.1
	github.com/cucumber/godog v0.11.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.7.0
)
