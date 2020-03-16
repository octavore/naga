# Naga

`naga` is a minimalistic service framework for Go. Build services by composing shared, testable, reusable modules.

```go
package main

import (
  "github.com/octavore/naga/service"
  "github.com/octavore/nagax/web"
)

func main() {
  service.Run(&web.Module{})
}
```

```shell
./simpleapp start
```

### Module

A **module** is a struct which implements the `service.Module` interface, which only has one method: `Init(*service.Config)`. The `Config` parameter allows configuration of lifecycle hooks and command-line options.

There are four hooks: `Setup`, `SetupTest`, `Start`, and `Stop`.

```go
package config

import (
  "encoding/json"
  "fmt"
  "io/ioutil"

  "github.com/octavore/naga/service"

// Module implements service.Module
type Module struct{ bytes []byte }

// Init configures a setup hook to read the config file
func (m *Module) Init(c *service.Config) {
  c.Setup = func() (err error) {
    m.bytes, err = ioutil.ReadFile("./config.json")
    return
  }
  c.Start = func() {
    fmt.Println(string(m.bytes))
  }
}

// Decode config data into i
func (m *Module) ReadConfig(i interface{}) error {
  return json.Unmarshal(m.bytes, i)
}
```

### Modules of Modules

Modules may depend on other modules, by adding them as a field. Dependencies are found via reflection and are recursively initialized.

Modules are singletons, so two modules which share the same dependency will be given the same instance of the dependency.

Lifecycle hooks in topological order (dependencies first when starting; dependencies last when stopping).

```go
type Module struct {
  Config     *config.Module
  HTTPConfig struct {
    Port string
  }
}

func (m *Module) Init(c *service.Config) {
  c.Setup = func() error {
    return m.Config.ReadConfig(&m.HTTPConfig)
  }
  c.Start = func() {
    addr := "localhost:" + m.HTTPConfig.Port
    err = http.ListenAndServe(addr, m)
    if err != nil {
      panic(err)
    }
  }
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
  rw.Write([]byte("hello world"))
}
```

### Commands

Your app's main function should initialize your app using `service.Run()`

```go
func main() {
  service.Run(&web.Module{})
}
```

This adds a `start` command which you can use to start your app. This will run all the `Start` hooks.

```shell
./myapp start
```

It's easy enough to register your own commands, which will run after the `Setup` hook.

```go
func (*AppModule) Init(c *service.Config) {
  c.AddCommand(&service.Command{
    Keyword: "domything",
    Run:     func(*service.CommandContext) {
      fmt.Println("hello world")
    },
  })
}
```

```shell
./myapp domything
```


### Environment

A Naga app will read the environment name from the `NAGA_ENV`. Supported environment names are `test`, `development`, and `production`.


