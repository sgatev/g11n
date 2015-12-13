package locale

// Loader represents a locale loader for a specific file format.
type Loader interface {

	// Load loads the locale from the file exposing a map of translated messages.
	Load(fileName string) map[string]string
}

var loaders = map[string]Loader{}

// Loaders returns all registered local loaders for different formats.
func Loaders() map[string]Loader {
	return loaders
}

// RegisterLoader registers a locale loader for specific format.
func RegisterLoader(format string, loader Loader) {
	loaders[format] = loader
}
