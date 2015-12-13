package locale

// Loader represents a locale loader for a specific file format.
type Loader interface {

	// Load loads the locale from the file exposing a map of translated messages.
	Load(fileName string) map[string]string
}

var loaders = map[string]Loader{}

// GetLoader returns the locale loader for a specific format.
func GetLoader(format string) (Loader, bool) {
	loader, ok := loaders[format]
	return loader, ok
}

// RegisterLoader registers a locale loader for specific format.
func RegisterLoader(format string, loader Loader) {
	loaders[format] = loader
}
