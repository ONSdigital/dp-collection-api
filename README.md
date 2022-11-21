
---

:warning: This repository will be archived in November 2022 as it is no longer in development. :warning:

---

dp-collection-api
================
An API for collection management

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable           | Default     | Description
| ------------------------------ | ----------- | -----------
| BIND_ADDR                      | :26000      | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s          | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL           | 30s         | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s         | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| MONGODB_COLLECTIONS_DATABASE   | collections | The MongoDB collections database
| MONGODB_COLLECTIONS_COLLECTION | collections | The MongoDB collections collection
| MONGODB_USERNAME               | test        | The MongoDB Username
| MONGODB_PASSWORD               | test        | The MongoDB Password
| MONGODB_CA_FILE_PATH           | file-path   | The MongoDB CA FilePath

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

