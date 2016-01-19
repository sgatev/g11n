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
//
// III. Format result
//
// The result of a message call could be further formatted by declaring a special
// result type that implements
//	G11nResult(formattedMessage string) string
// The format method is invoked after all parameters have been substituted in the message.
//
//	type SafeHTMLFormat string
//
//	func (shf SafeHTMLFormat) G11nResult(formattedMessage string) string {
//		r := strings.NewReplacer("<", `\<`, ">", `\>`, "/", `\/`)
//		return r.Replace(formattedMessage)
//	}
//
//	type M struct {
//		MyLittleSomething func() SafeHTMLFormat `default:"<message>Oops!</message>"`
//	}
package g11n
