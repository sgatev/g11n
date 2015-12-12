# g11n
[![Build Status](https://travis-ci.org/s2gatev/g11n.svg?branch=master)](https://travis-ci.org/s2gatev/g11n)
[![Coverage Status](https://coveralls.io/repos/s2gatev/g11n/badge.svg?branch=master&service=github)](https://coveralls.io/github/s2gatev/g11n?branch=master)

**g11n** */gopherization/* is an internationalization library inspired by [GWT](http://www.gwtproject.org/doc/latest/DevGuideI18nMessages.html) that offers:

* Statically-typed messages.
* Parameterized messages.
* Different options for defining messages' values.
* Extendable message formatting.

## Example

```go
package main

import (
	"fmt"

	"github.com/s2gatev/g11n"
)

type Messages struct {
	TheAnswer func(string, int) string `default:"The answer to %v is %v."`
}

func ExampleGopherization() {
	m := g11n.Init(&Messages{}).(*Messages)
	fmt.Print(m.TheAnswer("everything", 42))

	// Output:
	// The answer to everything is 42.
}
```

## License

The **g11n** library is licensed under the [MIT License](LICENSE).
