package apollo

import (
	"reflect"
)

//OptionsList - lista de options
type OptionsList []Options

// Options - opções para a injeção da dependência
type Options struct {
	key      key
	task     func(objectKey key, digo *Apollo)
	priority int
}

func (ol OptionsList) Len() int {
	return len(ol)
}

func (ol OptionsList) Less(i, j int) bool {
	return ol[i].priority < ol[j].priority
}

func (ol OptionsList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
}

// Singleton - determina se a dependência é um singleton, o padrão é true.
func Singleton(isSingleton bool) Options {
	return Options{
		task: func(objectKey key, digo *Apollo) {
			if !isSingleton {
				digo.cache.notSingletons[objectKey] = !isSingleton
			}
		},
		priority: 1,
	}
}

// As - especificar se construtor retorna uma interface
func As(typei interface{}) Options {
	return Options{
		task: func(objectKey key, digo *Apollo) {
			object, ok := digo.cache.components[objectKey]

			if !ok {
				object = digo.Fetch(objectKey.typed)
			}

			ikey := key{reflect.TypeOf(typei).Elem()}

			if _, ok := digo.cache.components[ikey]; !ok {
				digo.cache.components[ikey] = object
			}
		},
		priority: 2,
	}
}

// Qualifier - determina o tipo específico da dependência injetada como interface.
func Qualifier(typeObject, typeInterface interface{}) Options {
	return Options{
		task: func(objectKey key, digo *Apollo) {

			qualifierObjectKey := key{reflect.TypeOf(typeObject).Elem()}
			qualificerInterfaceKey := key{reflect.TypeOf(typeInterface)}

			qualifierObject := digo.Fetch(qualifierObjectKey.typed)
			qualifierObjectValue := reflect.ValueOf(qualifierObject)

			constructor := digo.cache.constructors[objectKey]
			argsOldConstructor := digo.getArgs(constructor)

			args := make([]reflect.Value, len(argsOldConstructor))

			for i := 0; i < len(argsOldConstructor); i++ {
				objectInterfaceType := qualificerInterfaceKey.typed.Elem()

				if argsOldConstructor[i].Implements(objectInterfaceType) {
					args[i] = qualifierObjectValue
				} else {
					args[i] = reflect.ValueOf(digo.Fetch(argsOldConstructor[i]))
				}
			}

			newConstructor := reflect.MakeFunc(reflect.TypeOf(constructor), func(arguments []reflect.Value) []reflect.Value {
				return reflect.ValueOf(constructor).Call(args)
			})

			digo.cache.constructors[objectKey] = newConstructor.Interface()

		},
		priority: 3,
	}
}
