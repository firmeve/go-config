# Point syntax configuration

[![Build Status](https://travis-ci.com/firmeve/go-config.svg?branch=master)](https://travis-ci.com/firmeve/go-config)
[![Coverage Status](https://coveralls.io/repos/github/firmeve/go-config/badge.svg?branch=master)](https://coveralls.io/github/firmeve/go-config?branch=master)
[![GitHub license](https://img.shields.io/github/license/firmeve/go-config.svg)](https://github.com/firmeve/go-config/blob/master/LICENSE)

Configuration package that supports multi-level dot syntax read write

## Basic usage

### Install
```go
go get -u github.com/firmeve/go-config
```

### Instantiation
```go
config, err := NewConfig(directory)
if err != nil {
    panic(err.Error())
}
```
> Note: `config` is a singleton object

### Set
// Set or add a value for the key
```go
err := config.Set("app.key", "value")
if err != nil {
    panic(err.Error())
}
```

// Set or add a value for a multi-section key
```go
err := config.Set("app.section.key", "value")
if err != nil {
    panic(err.Error())
}
```

### Get

```go
value, err := config.Get("app.section.key")
if err != nil {
    panic(err.Error())
}

fmt.Println("%#v",value)
fmt.Println(value.(*ini.Key).Value())

```
- If `key` is the file name, return a `*ini.File` object, such as: `Get("app")`
- If `key` has only one dot separation, it returns the `key` value specified by the default `section`, such as: `Get("app.t")`
- If `key` is a `section`, return a `*ini.Section` object, such as: `Get("app.section.sub")`

### GetDefault
Same as `Get`, but if the value accessed does not exist, the default value will be used instead of `err` to mask the error output.
```go
value := config.GetDefault("app.section.key", "default")
fmt.Println("%s",value)
```

### All
Get all configurations, return `map[string]*ini.File`
```
configs := config.All()
```

## Thanks
Development package based on [go-ini/ini](https://github.com/go-ini/ini)