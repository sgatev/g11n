package g11n

import (
	"fmt"
	"reflect"

	"github.com/s2gatev/g11n/locale"
)

const (
	embeddedMessageTag = "default"

	wrongResultsCountMessage = "Wrong number of results in a g11n message. Expected 1, got %v."
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

// materializeValue extracts the data from a reflected value and returns it.
func materializeValue(value reflect.Value) interface{} {
	i := value.Interface()

	if paramFormatter, ok := i.(paramFormatter); ok {
		return paramFormatter.G11nParam()
	}

	return i
}

// messageHandler creates a handler formats a message based on provided parameters.
func messageHandler(messagePattern string, resultType reflect.Type) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		var materializedArgs []interface{}
		for _, arg := range args {
			materializedArgs = append(materializedArgs, materializeValue(arg))
		}

		message := fmt.Sprintf(messagePattern, materializedArgs...)

		messageValue := reflect.ValueOf(message)
		resultValue := reflect.New(resultType).Elem()

		if resultFormatter, ok := resultValue.Interface().(resultFormatter); ok {
			modified := resultFormatter.G11nResult(message)
			modifiedValue := reflect.ValueOf(modified)
			resultValue.Set(modifiedValue.Convert(resultType))
		} else {
			resultValue.Set(messageValue)
		}

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
func (mf *MessageFactory) LoadLocale(format, localeName, fileName string) {
	loaders := locale.Loaders()
	if loader, ok := loaders[format]; ok {
		mf.locales[localeName] = loader.Load(fileName)
	} else {
		panic("Unknown format '" + format + "'.")
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
		messagePattern := field.Tag.Get(embeddedMessageTag)

		if locale, ok := mf.locales[mf.activeLocale]; ok {
			messageKey := fmt.Sprintf("%v.%v",
				concreteType.Name(),
				field.Name)
			if message, ok := locale[messageKey]; ok {
				messagePattern = message
			}
		}

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
