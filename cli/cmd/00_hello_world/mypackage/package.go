package mypackage

import "fmt"

func PublicFunction() {
	fmt.Println("Hello from the public function!")
}

func privateFunction() {
	fmt.Println("Hello from the private function!")
}

type Counter struct {
	count int
	Name  string
}

func NewCounter() *Counter {
	counter := &Counter{
		Name: "My Counter",
	}
	counter.setCount(0)
	return counter
}

func (c *Counter) setCount(count int) {
	c.count = count
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) GetCount() int {
	fmt.Printf("Counter: Name: %s, Value: %d\n", c.Name, c.count)
	return c.count
}
