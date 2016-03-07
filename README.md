# g11n

[![Join the chat at https://gitter.im/s2gatev/g11n](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/s2gatev/g11n?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/s2gatev/g11n.svg?branch=master)](https://travis-ci.org/s2gatev/g11n)
[![Coverage Status](https://coveralls.io/repos/s2gatev/g11n/badge.svg?branch=master&service=github)](https://coveralls.io/github/s2gatev/g11n?branch=master)
[![Go Report Card](http://goreportcard.com/badge/s2gatev/g11n)](http://goreportcard.com/report/s2gatev/g11n)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/s2gatev/g11n)
[![MIT License](http://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**g11n** */gopherization/* is an internationalization library inspired by [GWT](http://www.gwtproject.org/doc/latest/DevGuideI18nMessages.html) that offers:

* **Statically-typed** message keys.
* **Parameterized** messages.
* **Extendable** message formatting.
* **Custom** localization **file format**.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/s2gatev/g11n"
	locale "github.com/s2gatev/g11n/http"
)

type Messages struct {
	Hello func(string) string `default:"Hi %v!"`
}

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// Create messages factory.
		factory := g11n.New()

		// Initialize messages value.
		var m Messages
		factory.Init(&m)

		// Set messages locale.
		locale.SetLocale(factory, r)

		fmt.Fprintf(w, m.Hello("World"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
