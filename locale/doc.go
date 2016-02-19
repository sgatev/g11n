// Package locale presents localization loaders for the g11n internationalization library.
//
//
// I. Creating a locale loader
//
// Every locale loader (built-in or custom) should implement the Loader interface and register itself in the loaders registry under some specific format name using RegisterLoader.
//
//	func init() {
//		RegisterLoader("custom", customLoader{})
//	}
//
//
// II. Retrieving a locale loader
//
// A locale loader could be retrieved from the loaders registry using GetLoader.
//
//	loader, ok := GetLoader("custom") (Loader, bool)
//
//
// III. Built-in locale loaders
//
// g11n comes with two built-in locale loaders - "json" and "yaml".
package locale
