# Apollo: Dependency Injection for Go

[![Build Status](https://travis-ci.com/m2msolutionssa/apollo.svg?branch=master)]
[![Build Status](https://travis-ci.com/google/wire.svg?branch=master)]
[![godoc](https://godoc.org/github.com/m2msolutionssa/apollo?status.svg)]
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

