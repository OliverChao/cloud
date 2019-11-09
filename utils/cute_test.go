package utils

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string `tomap:"name"`
	Age  int    `tomap:"age"`
}

func TestStruct2Map(t *testing.T) {
	p := &Person{
		Name: "oliver",
		Age:  20,
	}

	struct2Map := Struct2Map(p)
	fmt.Println(struct2Map)
}
