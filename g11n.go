package g11n

import (
	"fmt"
	"reflect"
)

const (
	embeddedMessageTag = "embed"
)

// materializeValue extracts the data from a reflected value and returns it.
func materializeValue(value reflect.Value) interface{} {
	return value.Interface()
}

// messageHandler creates a handler formats a message based on provided parameters.
func messageHandler(messagePattern string) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		var materializedArgs []interface{}
		for _, arg := range args {
			materializedArgs = append(materializedArgs, materializeValue(arg))
		}

		message := fmt.Sprintf(messagePattern, materializedArgs...)

		messageValue := reflect.ValueOf(message)
		return []reflect.Value{messageValue}
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
		messageProxyFunc := reflect.MakeFunc(
			field.Type, messageHandler(messagePattern))
		instanceField.Set(messageProxyFunc)
	}

	return structPtr
}
