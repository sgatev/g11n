// Package g11n is an internationalization library that offers:
//
//	* Statically-typed message keys.
//	* Parameterized messages.
//	* Extendable message formatting.
//	* Custom localization file format.
//
//
// I. Initialization
//
// Create a new instance of g11n. Each instance handles messages and locales separately.
//
//	G = g11n.New()
//
// Define a struct with messages.
//
//	type Messages struct {
//		TheAnswer func(string, int) string `default:"The answer to %v is %v."`
//	}
//
// Initialize an instance of the struct through the g11n object.
//
//	var M *Messages
//	G.Init(M)
//
// Invoke messages on that instance.
//
//	M.TheAnswer("everything", 42")
//
//
// II. Choosing locale
//
// Specify the locale for every message struct initialized by this g11n instance.
//
//	G.SetLocale("en")
//
package g11n
