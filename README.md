#### Hexa contain services kit.

__Note : Hexa is in progress and does not have any stable release.__

#### Requirements:
go : minimum version `1.13`

#### Install
```
go get github.com/Kamva/hexa
```

#### Available Services:
- Config: config service.
- Logger: logger service
- Translator: translator service.

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

- [ ] We should implement Distributed tracing by using zipkin and open tracing and also use it in gRPC,... . but I think this
should be in the service mesh not in the business logic.

- [ ] Implement a service (e.g `Ping`) which should implement by all of the Hexa services to check health of that service [**Accepted**].

#### Todo
- [ ] Where we are checking log level on logger initialization step?
- [ ] Change all of kamva packages from uppercase to lowercase.
- [ ] Write Tests
- [ ] Add badges to readme.
- [ ] CI

#### Client conventions:
- [ ] on occure error and want to say to user that "some error occured",  give the requestID to the user, then user can give it back to the support team, and support team can see logs relative to that request id.
