package g11n

import (
	"fmt"
	"reflect"
)

const (
	embeddedMessageTag = "default"
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

// Init initializes the message fields of a structure pointer.
func Init(structPtr interface{}) interface{} {
	instance := reflect.ValueOf(structPtr).Elem()

	concreteType := instance.Type()
	for i := 0; i < concreteType.NumField(); i++ {
		field := concreteType.Field(i)
		instanceField := instance.FieldByName(field.Name)
		messagePattern := field.Tag.Get(embeddedMessageTag)
		resultType := field.Type.Out(0)
		messageProxyFunc := reflect.MakeFunc(
			field.Type, messageHandler(messagePattern, resultType))
		instanceField.Set(messageProxyFunc)
	}

	return structPtr
}
