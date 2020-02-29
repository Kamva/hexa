#### kitty contain services kit.

#### Install
```
go get github.com/Kamva/kitty
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

config := kittyconfig.NewViperDriver(v)

// Use config service in app.
```

#### today todo: 
- [ ] On worker and jobs implementations get queue prefix 
(in both jobs and worker, to distinct microservices queues).

#### Todo:
- [ ] Collection presenter
- [ ] Pagination
- [ ] Write Tests
- [ ] Add more services like log,...
- [ ] Add badges to readme.
- [ ] CI 