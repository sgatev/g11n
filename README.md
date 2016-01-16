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

	"github.com/s2gatev/g11n"
)

var G = g11n.New()

type Messages struct {
	TheAnswer func(string, int) string `default:"The answer to %v is %v."`
}

func ExampleGopherization() {
	m := G.Init(&Messages{}).(*Messages)
	fmt.Print(m.TheAnswer("everything", 42))

	// Output:
	// The answer to everything is 42.
}
```
