# this
`this` is a golang library to provide dead simple BDD-style context to yours tests. It works entirely within the default golang test framework. The goal here is to provide better color around what is actually being tested.

## Installation
You shouldn't have to do anything other than a standard golang library installation:
`get get -u github.com/connerhansen/this`

## Usage
`this` is easy to use. Just create a normal test as you would for golang, but then put individual test units in `this.Should` blocks:

```golang
import (
  "github.com/connerhansen/this"
)

func TestSomeFunctionality(t *testing.T) {

  this.After(func() {
    // Do something after each test
  })

  this.Before(func() {
    // Do something before each test
  })

  // All of your test logic is contained in this should statement
  this.Should("Return true if some condition is met", t,
    func() {
      ...
    })

}
```

If you want to use this with something like Gomega (strongly suggested), `this` also provides a default fail handler than you can register from within an `init()` func in your tests:

```golang
import (
  "github.com/connerhansen/this"
	. "github.com/onsi/gomega"
)

func init() {
	RegisterFailHandler(this.GomegaFailHandler)
}

func TestSomeFunctionality(t *testing.T) {

  // All of your test logic is contained in this should statement
  this.Should("Return true if some condition is met", t,
    func() {
      ...
      Expect(someVal).To(Equal(anotherVal))
    })

}
```
