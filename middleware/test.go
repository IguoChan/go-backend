package middleware

import "fmt"

type A struct {
	M map[string]interface{}
}

func main() {
	a := &A{}
	b, ok := a.M["nihao"].(int)
	if !ok {
		b = 10
	}
	fmt.Println(b)
}
