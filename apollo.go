package apollo

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
)

type Apollo struct {
	cache
}

type key struct {
	typed reflect.Type
}

type constructor interface{}
type component interface{}

type cache struct {
	components    map[key]component
	constructors  map[key]constructor
	notSingletons map[key]bool
	options       OptionsList
}

// New - inicializa o container de injeção de dependencias
func New() Apollo {
	return Apollo{
		cache{
			components:    map[key]component{},
			constructors:  map[key]constructor{},
			options:       OptionsList{},
			notSingletons: map[key]bool{},
		},
	}
}

func (apollo *Apollo) proccessOptions() {

	sort.Sort(apollo.cache.options)

	for i := 0; i < len(apollo.cache.options); i++ {
		op := apollo.cache.options[i]
		task := op.task
		task(op.key, apollo)
	}
}

func (apollo *Apollo) invokeConstructor(constructor constructor, args []reflect.Type) (returns []reflect.Value) {

	dependencies := make([]reflect.Value, len(args))

	for i := 0; i < len(dependencies); i++ {
		dependencie := apollo.Fetch(args[i])
		dependencies[i] = reflect.ValueOf(dependencie)
	}

	constructorValue := reflect.ValueOf(constructor)

	returns = constructorValue.Call(dependencies)
	return
}

func (apollo *Apollo) getFunctionName(valueFunc reflect.Value) string {
	return runtime.FuncForPC(valueFunc.Pointer()).Name()
}

func (apollo *Apollo) validateConstructor(valueFunc reflect.Value) error {

	typed := valueFunc.Type()

	if typed.Kind() != reflect.Func {
		return fmt.Errorf("O construtor %s deve ser uma função", typed.Kind().String())
	}

	if typed.NumOut() > 2 || typed.NumOut() < 1 {
		return fmt.Errorf("O numero de elementos retornados no construtor %s deve ser no mínimo 1 e no máximo 2", apollo.getFunctionName(valueFunc))
	}

	if typed.NumOut() == 2 {

		errorInterface := reflect.TypeOf((*error)(nil)).Elem()
		secondTypeReturn := typed.Out(1)

		if !secondTypeReturn.Implements(errorInterface) {
			return fmt.Errorf("O segundo elemento retornado no construtor %s deve ser um error", apollo.getFunctionName(valueFunc))
		}
	}

	return nil
}

func (apollo *Apollo) getArgs(constructor constructor) []reflect.Type {

	typec := reflect.TypeOf(constructor)

	deps := make([]reflect.Type, typec.NumIn())

	for i := 0; i < len(deps); i++ {
		deps[i] = typec.In(i)
	}

	return deps
}

// Register - registra dependências no container
func (apollo *Apollo) Register(constructor constructor, options ...Options) {

	valued := reflect.Indirect(reflect.ValueOf(constructor))

	err := apollo.validateConstructor(valued)

	if err != nil {
		panic(err)
	}

	k := key{valued.Type().Out(0)}

	optionsRegistered := []Options{}

	for _, op := range options {
		op.key = k
		optionsRegistered = append(optionsRegistered, op)
	}

	apollo.cache.constructors[k] = constructor
	apollo.cache.options = append(apollo.cache.options, optionsRegistered[0:len(options)]...)

}

// Fetch - constroi dependência solicitada
func (apollo *Apollo) Fetch(typed reflect.Type) interface{} {
	key := key{typed}

	if object, found := apollo.cache.components[key]; found {
		return object
	}

	constructor := apollo.cache.constructors[key]
	args := apollo.getArgs(constructor)

	returns := apollo.invokeConstructor(constructor, args)

	if len(returns) == 2 && !returns[1].IsNil() {
		err := returns[1].Interface().(error)
		panic(err)
	}

	object := returns[0].Interface()

	isNotSingleton := apollo.cache.notSingletons[key]

	if !isNotSingleton {
		apollo.cache.components[key] = object
	}

	return object
}

// Init - inicializa o container de injeção de dependências
func (apollo *Apollo) Init(initFunc constructor) {

	apollo.proccessOptions()
	args := apollo.getArgs(initFunc)
	apollo.invokeConstructor(initFunc, args)
}
