package g11n

import (
	"fmt"
	"reflect"
	"sync"

	g11nLocale "github.com/s2gatev/g11n/locale"
)

const /* application constants */ (
	defaultMessageTag = "default"
)

const /* error message patterns */ (
	wrongResultsCountMessage = "Wrong number of results in a g11n message. Expected 1, got %v."
	unknownFormatMessage     = "Unknown locale format '%v'."
)

// Synchronizer synchronizes asynchronous tasks.
type Synchronizer struct {
	tasks     *sync.WaitGroup
	completed bool
}

// Await awaits the completion of the tasks.
func (s *Synchronizer) Await() {
	s.tasks.Wait()
}

// Completed returns whether the tasks are already completed.
func (s *Synchronizer) Completed() bool {
	return s.completed
}

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

// formatParam extracts the data from a reflected argument value and returns it.
func formatParam(value reflect.Value) interface{} {
	valueInterface := value.Interface()

	if paramFormatter, ok := valueInterface.(paramFormatter); ok {
		return paramFormatter.G11nParam()
	}

	return valueInterface
}

// messageHandler creates a handler that formats a message based on provided parameters.
func messageHandler(messagePattern string, resultType reflect.Type) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		// Format message parameters.
		var formattedParams []interface{}
		for _, arg := range args {
			formattedParams = append(formattedParams, formatParam(arg))
		}

		// Find the result message value.
		message := fmt.Sprintf(messagePattern, formattedParams...)
		messageValue := reflect.ValueOf(message)

		// Format message result.
		resultValue := reflect.New(resultType).Elem()
		if resultFormatter, ok := resultValue.Interface().(resultFormatter); ok {
			formattedResult := resultFormatter.G11nResult(message)
			messageValue = reflect.ValueOf(formattedResult).Convert(resultType)
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
	mf.initializeStruct(structPtr)

	return structPtr
}

// InitAsync initializes the message fields of a structure pointer asynchronously.
func (mf *MessageFactory) InitAsync(structPtr interface{}) (interface{}, *Synchronizer) {
	var initializers sync.WaitGroup
	synchronizer := &Synchronizer{tasks: &initializers}

	initializers.Add(1)
	go func() {
		mf.initializeStruct(structPtr)
		initializers.Done()
		synchronizer.completed = true
	}()

	return structPtr, synchronizer
}

// initializeStruct initializes the message fields of a struct pointer.
func (mf *MessageFactory) initializeStruct(structPtr interface{}) {
	instance := reflect.ValueOf(structPtr).Elem()
	concreteType := instance.Type()

	// Initialize each message func of the struct.
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

		// Create proxy function for handling the message.
		messageProxyFunc := reflect.MakeFunc(
			field.Type, messageHandler(messagePattern, resultType))

		instanceField.Set(messageProxyFunc)
	}
}
