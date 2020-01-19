// main.go

package main

import "fmt"

func main() {
	a := App{}
	a.Initialize()

	fmt.Println("Serving")
	a.Run(":8080")
}
