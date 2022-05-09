package dino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myInterface1 interface {
	Method1()
}

type myStruct1 struct {
	Foo int
}

func (s *myStruct1) Method1() {}

type myStruct2 struct{}

func TestAddFailsServiceType(t *testing.T) {
	var err error

	err = Add[myStruct1, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = Add[int, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = Add[string, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = Add[*myInterface1, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})
	assert.Contains(t, err.Error(), "*dino.myInterface1")
}

func TestAddFailsNotImplements(t *testing.T) {
	err := Add[myInterface1, myStruct2](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &NotImplementsError{})
	assert.Regexp(t, "myInterface1.*myStruct2", err.Error())
}

func TestAddFailsImplNotStruct(t *testing.T) {
	err := Add[*myStruct1, *myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &ImplNotStructError{})
	assert.Contains(t, err.Error(), "*dino.myStruct1")
}

func TestAddFailsBadPointer(t *testing.T) {
	err := Add[*myStruct1, myStruct2](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &BadPointerError{})
	assert.Contains(t, err.Error(), " dino.myStruct1")
}

func TestAddSucceedsAndIsSingleton(t *testing.T) {
	c := &Container{}
	err := Add[myInterface1, myStruct1](c)
	assert.Nil(t, err)

	s, err := Get[myInterface1](c)
	assert.Nil(t, err)
	assert.IsType(t, &myStruct1{}, s)

	s.(*myStruct1).Foo = 4

	s2, err := Get[myInterface1](c)
	assert.Nil(t, err)
	assert.IsType(t, &myStruct1{}, s)
	assert.Same(t, s, s2)

	assert.Equal(t, 4, s2.(*myStruct1).Foo)
	s.(*myStruct1).Foo = 5
	assert.Equal(t, 5, s2.(*myStruct1).Foo)
}

func TestAddInstanceFailsServiceType(t *testing.T) {
	var err error

	err = AddInstance[myStruct1](&Container{}, myStruct1{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = AddInstance[int](&Container{}, myStruct1{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = AddInstance[string](&Container{}, myStruct1{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = AddInstance[*myInterface1](&Container{}, myStruct1{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})
	assert.Contains(t, err.Error(), "*dino.myInterface1")
}

func TestAddInstanceFailsNotImplements(t *testing.T) {
	err := AddInstance[myInterface1](&Container{}, myStruct2{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &NotImplementsError{})
	assert.Regexp(t, "myInterface1.*myStruct2", err.Error())
}

func TestAddInstanceFailsBadPointer(t *testing.T) {
	err := AddInstance[*myStruct1](&Container{}, &myStruct2{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &BadPointerError{})
	assert.Contains(t, err.Error(), " dino.myStruct1")
}

func TestAddInstanceStructPtrSucceeds(t *testing.T) {
	c := &Container{}
	instance := &myStruct1{Foo: 5}
	err := AddInstance[*myStruct1](c, instance)
	assert.Nil(t, err)
	svc, err := Get[*myStruct1](c)
	assert.Nil(t, err)
	instance.Foo = 6
	assert.Equal(t, 6, svc.Foo)
}

func TestAddTransientFailsServiceType(t *testing.T) {
	var err error

	err = AddTransient[myStruct1, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = AddTransient[int, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = AddTransient[string, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})

	err = AddTransient[*myInterface1, myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &InvalidServiceTypeError{})
	assert.Contains(t, err.Error(), "*dino.myInterface1")
}

func TestAddTransientFailsNotImplements(t *testing.T) {
	err := AddTransient[myInterface1, myStruct2](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &NotImplementsError{})
	assert.Regexp(t, "myInterface1.*myStruct2", err.Error())
}

func TestAddTransientFailsImplNotStruct(t *testing.T) {
	err := AddTransient[*myStruct1, *myStruct1](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &ImplNotStructError{})
	assert.Contains(t, err.Error(), "*dino.myStruct1")
}

func TestAddTransientFailsBadPointer(t *testing.T) {
	err := AddTransient[*myStruct1, myStruct2](&Container{})
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &BadPointerError{})
	assert.Contains(t, err.Error(), " dino.myStruct1")
}

func TestAddTransientErrorsCyclicDependency(t *testing.T) {
	type (
		MyIf interface{}
		X    struct {
			ZDep MyIf
		}
		Y struct {
			XDep *X
		}
		Z struct {
			YDep *Y
		}
	)

	c := &Container{}
	assert.Nil(t, AddTransient[*X, X](c))
	assert.Nil(t, AddTransient[*Y, Y](c))
	assert.Nil(t, AddTransient[MyIf, Z](c))

	_, err := Get[MyIf](c)
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &CyclicDependencyError{})
	assert.Regexp(t, "MYIF .transient.$", err.Error())
}
