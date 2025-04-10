package main

import (
	"errors"
	"fmt"
)

func add(a int, b int) int {
	return a + b
}

func errorIfEmpty(value string) (string, error) {
	if value == "" {
		return "", errors.New("value is empty")
	}
	return value, nil
}

func main() {
	// 1. print hello world
	fmt.Println("Hello, World!")

	// // 2. declare variable with default value
	// var a string
	// fmt.Println(a)
	// fmt.Println(a == "")

	// // 3. declare variable with initial value
	// var b string = "initial"
	// fmt.Println(b)

	// // 4. declare variable with type inference
	// var c = "type inference"
	// fmt.Println(c)

	// // 5. declare variable with short declaration
	// d := "short"
	// fmt.Println(d)

	// // 6. call function and print result
	// result := add(1, 2)
	// fmt.Println(result)

	// // 7. call function and handle error
	// _, err := errorIfEmpty("")
	// if err != nil {
	// 	// handle error or return error
	// 	fmt.Println(err)
	// }

	// // 8. string interpolation
	// fmt.Printf("Hello, %s\n", "World")

	// // 9. package import
	// mypackage.PublicFunction()
	// // mypackage.privateFunction() // This will not work because it's private

	// counter := mypackage.NewCounter()
	// counter.Increment()
	// fmt.Println(counter.GetCount())
	// counter.Increment()
	// fmt.Println(counter.GetCount())
}
