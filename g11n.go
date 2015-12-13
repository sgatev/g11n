package g11n

import (
	"fmt"
	"reflect"

	g11nLocale "github.com/s2gatev/g11n/locale"
)

const /* application constants */ (
	defaultMessageTag = "default"
)

const /* error message patterns */ (
	wrongResultsCountMessage = "Wrong number of results in a g11n message. Expected 1, got %v."
	unknownFormatMessage     = "Unknown locale format '%v'."
)

// paramFormatter represents a type that supports custom formatting
// when it is used as parameter in a call to a g11n message.
type paramFormatter interface {

	// G11nParam formats a type in a specific way when passed to a g11n message.
	G11nParam() string
}

// resultFormatter represents a type that supports custom formatting
// when it is returned from a call to a g11n message.
type resultFormatter interface {

	// G11nResult accepts a formatted g11n message and modifies it before returning.
	G11nResult(formattedMessage string) string
}

// formatArg extracts the data from a reflected argument value and returns it.
func formatArg(value reflect.Value) interface{} {
	valueInterface := value.Interface()

	if paramFormatter, ok := valueInterface.(paramFormatter); ok {
		return paramFormatter.G11nParam()
	}

	return valueInterface
}

// messageHandler creates a handler formats a message based on provided parameters.
func messageHandler(messagePattern string, resultType reflect.Type) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		resultValue := reflect.New(resultType).Elem()

		// Format message arguments.
		var formattedArgs []interface{}
		for _, arg := range args {
			formattedArgs = append(formattedArgs, formatArg(arg))
		}

		// Find the result message value.
		message := fmt.Sprintf(messagePattern, formattedArgs...)
		messageValue := reflect.ValueOf(message)
		if resultFormatter, ok := resultValue.Interface().(resultFormatter); ok {
			modified := resultFormatter.G11nResult(message)
			modifiedValue := reflect.ValueOf(modified)
			messageValue = modifiedValue.Convert(resultType)
		}

		resultValue.Set(messageValue)

		return []reflect.Value{resultValue}
	}
}

// MessageFactory initializes message structs and provides language
// translations to messages.
type MessageFactory struct {
	activeLocale string
	locales      map[string]map[string]string
}

// New returns a fresh G11n message factory.
func New() *MessageFactory {
	return &MessageFactory{
		locales: map[string]map[string]string{},
	}
}

// LoadLocale loads the content of a locale file in the specified format.
func (mf *MessageFactory) LoadLocale(format, locale, fileName string) {
	if loader, ok := g11nLocale.GetLoader(format); ok {
		mf.locales[locale] = loader.Load(fileName)
	} else {
		panic(fmt.Sprintf(unknownFormatMessage, format))
	}
}

// SetLocale sets the currently active locale for the messages generated
// by this factory.
func (mf *MessageFactory) SetLocale(locale string) {
	mf.activeLocale = locale
}

// Init initializes the message fields of a structure pointer.
func (mf *MessageFactory) Init(structPtr interface{}) interface{} {
	instance := reflect.ValueOf(structPtr).Elem()

	concreteType := instance.Type()
	for i := 0; i < concreteType.NumField(); i++ {
		field := concreteType.Field(i)
		instanceField := instance.FieldByName(field.Name)

		// Extract default message.
		messagePattern := field.Tag.Get(defaultMessageTag)

		// Extract localized message.
		if locale, ok := mf.locales[mf.activeLocale]; ok {
			messageKey := fmt.Sprintf("%v.%v", concreteType.Name(), field.Name)
			if message, ok := locale[messageKey]; ok {
				messagePattern = message
			}
		}

		// Check if return type of the message func is correct.
		if field.Type.NumOut() != 1 {
			panic(fmt.Sprintf(wrongResultsCountMessage, field.Type.NumOut()))
		}

		resultType := field.Type.Out(0)
		messageProxyFunc := reflect.MakeFunc(
			field.Type, messageHandler(messagePattern, resultType))
		instanceField.Set(messageProxyFunc)
	}

	return structPtr
}
