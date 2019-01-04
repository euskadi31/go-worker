Go Worker [![Last release](https://img.shields.io/github/release/euskadi31/go-worker.svg)](https://github.com/euskadi31/go-worker/releases/latest) [![Documentation](https://godoc.org/github.com/euskadi31/go-worker?status.svg)](https://godoc.org/github.com/euskadi31/go-worker)
=========

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-worker)](https://goreportcard.com/report/github.com/euskadi31/go-worker)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/euskadi31/go-worker/master.svg)](https://travis-ci.org/euskadi31/go-worker) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-worker/master.svg)](https://coveralls.io/github/euskadi31/go-worker?branch=master) |

Pool Worker fo Golang

## Example

```go
import (
    "fmt"

    "github.com/euskadi31/go-worker"
)

pool := worker.New(10, 100, func(payload interface{}) {
    val := payload.(int)

    fmt.Printf("Val: %d", val * val)
})

pool.Start()

pool.Enqueue(1)

pool.Close()

```

## License

go-worker is licensed under [the MIT license](LICENSE.md).
