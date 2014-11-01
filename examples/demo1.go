package main

import (
	"fmt"
	"github.com/ukautz/mappath"
)

func main() {
	c := mappath.NewMapPath(map[string]interface{}{
		"foo": "bar",
		"baz": map[string]interface{}{
			"foo": "bar",
			"baz": map[string]interface{}{
				"foo": "bar",
				"baz": "bar",
			},
		},
	})
	v, _ := c.Get("baz/baz/baz")

	// catch path not found
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Got whole: %+v\n", v)
	}
}
