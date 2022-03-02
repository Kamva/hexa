#### Hexa is a microservice SDK.

__Note : Hexa is in progress, it does not have stable version until now, so use it at your own risk.__

#### Requirements:

go : minimum version `1.13`

#### Install

```
go get github.com/kamva/hexa
```

#### Available Services:

- Config: config service.
- Logger: logger service and exception tracker (using Sentry)
- Translator: translator service.
- Distributed Lock Manager: At the moment we support __MongoDB__ as driver, but if needed we will Add __Etcd__ and  __
  Redis red locks__.
- [Event Server](http://github.com/kamva/hexa-event):  available drivers are:
    - Kafka (and __kafka outbox__)
    - pulsar
    - NATS streaming
- [Hexa Arranger to manage your workflow using Temporal](https://github.com/Kamva/hexa-arranger)
- [Tools and interceptors for gRPC](https://github.com/Kamva/hexa-rpc)
- [Hexa Tuner](https://github.com/Kamva/hexa-tuner) boot all Hexa services without headache
- [Hexa Sendo](https://github.com/Kamva/hexa-sendo) Send SMS, email and push notifications(todo) in your
  microservices.
- [Hexa Jobs](https://github.com/Kamva/hexa-job) push your jobs and use queues to handle jobs.
- [Hexa tools for routing and web server](https://github.com/Kamva/hexa-echo)

#### How to use:

example:

```go
// Assume we want to use viper as config service.
v := viper.New()

// tune your viper.
config := hexaconfig.NewViperDriver(v)

// Use config service in app.
```

#### Proposal

- [ ] Replace http status code with gRPC status code in our errors (also maybe in replies).

- [ ] Implement Hexa Context propagators for distributed tracing.

- [x] Implement a service (e.g `Health`) which should implement by all Hexa services to check health of that
  service [**Accepted**][Implemented].


### Notes
- Run `go generate ./db/...` or `go generate ./...` to generate some files after your changes in the DB models.

#### Todo

- [x] Where are we checking log level on logger initialization step?
- [x] Change all of kamva packages from uppercase to lowercase.
- [ ] Write Tests
- [ ] Add badges to readme.
- [ ] CI for tests.
