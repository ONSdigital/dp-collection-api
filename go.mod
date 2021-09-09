module github.com/ONSdigital/dp-collection-api

go 1.16

replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible

require (
	github.com/ONSdigital/dp-component-test v0.4.2-0.20210908155047-7f9ad0ffda5d
	github.com/ONSdigital/dp-healthcheck v1.1.0
	github.com/ONSdigital/dp-mongodb/v2 v2.2.0
	github.com/ONSdigital/dp-net/v2 v2.2.0-beta
	github.com/ONSdigital/log.go/v2 v2.0.6
	github.com/benweissmann/memongo v0.1.1
	github.com/cucumber/godog v0.11.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/afero v1.6.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	go.mongodb.org/mongo-driver v1.7.1
)
