#### kitty contain services kit.

#### Install
```
go get github.com/Kamva/kitty
```

#### Available Services:
- Config: to use as app config service.
- Translator: to use as translator.

#### How to use:
example:
```go
// Assume we want to use viper as config service.
v := viper.New()

// tune your viper.

config := kittyconfig.NewViperDriver(v)

// Use config service in app.
```

#### Todo:
- [ ] Write Tests
- [ ] Add more services like log,...
- [ ] Add badges to readme.
- [ ] CI 