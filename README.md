#### Hexa is a microservice SDK.

__Note : Hexa is in progress, it does not have stable version until now, so use it at your risk.__

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

- [ ] We should implement Distributed tracing by using zipkin and open tracing and also use it in gRPC,... . but I think
  this should be in the service mesh not in the business logic.

- [ ] Implement a service (e.g `Ping`) which should implement by all of the Hexa services to check health of that
  service [**Accepted**].

#### Todo

- [ ] Where we are checking log level on logger initialization step?
- [ ] Change all of kamva packages from uppercase to lowercase.
- [ ] Write Tests
- [ ] Add badges to readme.
- [ ] CI

#### Client conventions:

- [ ] on occure error and want to say to user that "some error occured", give the requestID to the user, then user can
  give it back to the support team, and support team can see logs relative to that request id.
