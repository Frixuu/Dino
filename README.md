# Dino

[![Go Report Card](https://goreportcard.com/badge/github.com/frixuu/dino)](https://goreportcard.com/report/github.com/frixuu/dino)
![GitHub](https://img.shields.io/github/license/frixuu/dino)
![Lines of code](https://img.shields.io/tokei/lines/github/frixuu/dino)

A dependency injection container for Go 1.18.

## Example

```golang
package main

import (
    "os"
    
    "github.com/frixuu/dino"
    "go.uber.org/zap"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type AccountIdCache interface{}
type AccountIdCacheImpl struct{}

type AccountController struct {
    Cache  AccountIdCache
    DB     *gorm.DB       `dino:"named:accounts"`
    Logger *zap.Logger
}

func (c *AccountController) DoWork() {
    c.Logger.Info("Hello, world!")
}

func main() {

    // Create the container
    c := &dino.Container{}

    // Register a singleton.
    // It will be created once and persist for the whole lifetime of the container
    dino.Add[AccountIdCache, AccountIdCacheImpl](c)

    // Register a transient.
    // It will be recreated each time it gets requested from the container
    dino.AddTransient[*AccountController, AccountController](c)

    // If you have some existing objects, you can register them as instances
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    dino.AddInstance[*zap.Logger](c, logger)

    // If you want to use the same types in different contexts, use named bindings
    db, _ := gorm.Open(postgres.Open(os.Getenv("DSN_ACCOUNTS")), &gorm.Config{})
    dino.AddInstanceNamed[*gorm.DB](c, "accounts", db)

    // Request a service from the container
    controller, _ := dino.Get[*AccountController](c)
    controller.DoWork()
}
```

## Credits

This project is influenced by [zekroTJA](https://github.com/zekroTJA/di)'s prior work, [MIT-licensed](https://github.com/zekroTJA/di/blob/390e0870d20ed665f4773b3c86ee0ee80eeeb352/LICENSE).
