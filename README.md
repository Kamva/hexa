#### Hexa contain services kit.

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

#### today todo: 
- [ ] On worker and jobs implementations get queue prefix 
(in both jobs and worker, to distinct microservices queues).

#### Todo:
- [ ] Collection presenter
- [ ] Write Tests
- [ ] Add more services like log,...
- [ ] Add badges to readme.
- [ ] CI 

#### Client conventions:
- [ ] on occure error and want to say to user that "some error occured",  give the requestID to the user, then user can give it back to the support team, and support team can see logs relative to that request id.