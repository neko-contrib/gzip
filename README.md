#gzip
[![wercker status](https://app.wercker.com/status/c14315e0509943fa722a01b1f06b8ebc/s "wercker status")](https://app.wercker.com/project/bykey/c14315e0509943fa722a01b1f06b8ebc)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/neko-contrib/gzip)

Gzip middleware for [Neko](https://github.com/rocwong/neko)

## Usage

~~~ go
package main
import (
  "github.com/rocwong/neko"
  "github.com/neko-contrib/gzip"
)

func main() {
  m := neko.New()
  m.Use(gzip.Gzip(gzip.DefaultCompression))
  m.Run(":3000")
}

~~~

Make sure to include the Gzip middleware above other middleware that alter the response body (like the render middleware).

## Authors
* [Jeremy Saenz](http://github.com/codegangsta)
* [Shane Logsdon](http://github.com/slogsdon)