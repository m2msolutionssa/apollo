# Apollo: Dependency Injection for Go

[![Build Status](https://travis-ci.com/m2msolutionssa/apollo.svg?branch=master)](https://travis-ci.com/m2msolutionssa/apollo)
[![godoc](https://godoc.org/github.com/m2msolutionssa/apollo?status.svg)](https://godoc.org/github.com/m2msolutionssa/apollo)
[![Coverage](https://codecov.io/gh/m2msolutionssa/apollo/branch/master/graph/badge.svg)](https://codecov.io/gh/m2msolutionssa/apollo)

Apollo is a reflection-based dependency injection library designed for Go projects.
Dependencies between components are represented in Apollo as constructor function 
parameters, encouraging explicit initialization rather than global variables.
Some features of Apollo were based on the Spring Framework, such as non-singleton 
dependency injection and interface-based dependency injection with the ability to specify the implementation.

## Installing

Install Apollo by running:

```shell
go get github.com/m2msolutionssa/apollo
```
and ensuring that `$GOPATH/bin` is added to your `$PATH`.

## Project status

Apollo is currently in *beta*. During the beta period, we encourage you to use Apollo and provide feedback. 
We will focus on improving and evolving the library as the needs of the community.

## Usage

The Apollo dependency injection container can be created as follows.

```go
container := apollo.New()
```
Apollo is based on building dependencies through constructor functions. Below is an example of creation.

```go
type Person struct {
  name string 
}

func NewPerson()Person{
  return Person{"Bob"}
}
```
Constructor functions must be registered via the Register method of the container.

```go
container.Register(NewPerson)
```
The container can be initialized using the Init method.

```go
container.Init(func(person Person){
  fmt.Println(person.name)
})
```
The Init method initializes the container and builds the dependency tree.

### Injecting Dependencies

Dependencies must be injected through constructor functions

```go
type Person struct {
  name string 
}

type Car struct{
  name string 
  owner Person
}

func NewPerson()Person{
  return Person{"Bob"}
}

func NewCar(person Person)Car{
  return Car{name:"Ferrari",owner:person}
}

container.Register(NewPerson)
container.Register(NewCar)

container.Init(func(car Car){
  fmt.Println(car.name,car.owner.name)
})

```
Optionally, constructor functions may return a second argument if an error has to be flagged in the dependency build

```go
func NewCar(person Person)(Car,error){
  if person.name != "Bob"{
    return Car{}, errors.New("The person's name must be Bob")
  }else{
    return Car{name:"Ferrari",owner:person}
  }
}
```
### Non singleton dependencies

By default, all injected dependencies are singleton, but you can change this behavior through Apollo options.

```go
type Person struct {
  name string 
}

type Car struct{
  name string 
  owner Person
}

type Airplane struct {
  owner Person
}

func NewPerson()Person{
  return Person{"Bob"}
}

func NewCar(person Person)Car{
  return Car{name:"Ferrari",owner:person}
}

func NewAirplane(person Person)Airplane{
  return Airplane{owner:person}
}

container.Register(NewPerson, apollo.Singleton(false))
container.Register(NewCar)
container.Register(NewAirplane)

container.Init(func(car Car, airplane Airplane){
  fmt.Println(&car.owner,&airplane.owner)
})

```
Different instances of Person are injected into Airplane and Car.

### Interface-based dependencies

Apollo allows components to request dependencies that implement a given interface.

```go
type Listener interface {
  observe() 
}

type Observer struct{}

type Subject struct {
  observer Observer 
}

func (o Observer) observe(){
  fmt.Println("watching ...")
}

func NewObserver()Observer{
  return Observer{}
}

func NewSubject(observer Listener)Subject{
  return Subject{observer:observer}
}

container := apollo.New() 
container.Register(NewObserver,apollo.As(new(Listener)))
container.Register(NewSubject)

container.Init(func(sub Subject){
  fmt.Println(sub.observer.observe())
})
  
```
The As function specifies that the default implementation for the Listener interface is the Observer structure.

### Multiple interface implementations

Apollo also allows you to specify alternative implementations for the same interface beyond the default implementation through the Qualifier function.

```go
type Listener interface {
  observe() 
}

type Observer struct{}

type Monitor struct{}

type Subject struct {
  observer Observer 
}

type Subject2 struct {
  observer Observer 
}

func (o Observer) observe(){
  fmt.Println("watching Observer ...")
}

func (o Monitor) observe(){
  fmt.Println("watching Monitor ...")
}

func NewObserver()Observer{
  return Observer{}
}

func NewMonitor()Monitor{
  return Monitor{}
}

func NewSubject(observer Listener)Subject{
  return Subject{observer:observer}
}

func NewSubject2(observer Listener)Subject2{
  return Subject2{observer:observer}
}

container := apollo.New() 
container.Register(NewObserver,apollo.As(new(Listener)))
container.Register(NewSubject)
container.Register(NewSubject2,apollo.Qualifier(new(Monitor),new(Listener)))
container.Register(NewMonitor)

container.Init(func(sub Subject, sub2 Subject2){
  fmt.Println(sub.observer.observe())
  fmt.Println(sub2.observer.observe())
})
  
```


