package main

import (
	"fmt"
	"github.com/ukautz/mappath"
	//"reflect"
)

// Working with nested arrays
func main() {
	c := mappath.NewMapPath(map[string]interface{}{
		"maps": []map[string]interface{}{
			map[string]interface{}{
				"foo": "bar1",
			},
			map[string]interface{}{
				"foo": "bar2",
			},
			map[string]interface{}{
				"foo": "bar3",
			},
			map[string]interface{}{
				"foo": "bar4",
			},
		},
		"ints":    []int{1, 2, 3, 4},
		"floats":  []float64{1.01, 2.02, 3.03, 4.04},
		"strings": []string{"one", "two", "three", "four"},
	})

	fmt.Println("Reading int array")
	ints, err := c.Ints("ints")
	if err != nil {
		fmt.Printf("  Error getting int array: %+v\n", err)
	} else {
		for i, val := range ints {
			fmt.Printf("  Int %d is: %d\n", i, val)
		}
	}

	fmt.Println("\nReading float array")
	floats, err := c.Floats("floats")
	if err != nil {
		fmt.Printf("  Error getting float array: %+v\n", err)
	} else {
		for i, val := range floats {
			fmt.Printf("  Float %d is: %f\n", i, val)
		}
	}

	fmt.Println("\nReading string array")
	strings, err := c.Strings("strings")
	if err != nil {
		fmt.Printf("  Error getting string array: %+v\n", err)
	} else {
		for i, val := range strings {
			fmt.Printf("  String %d is: %s\n", i, val)
		}
	}

	fmt.Println("\nReading map array")
	maps, err := c.Maps("maps")
	if err != nil {
		fmt.Printf("  Error getting map array: %+v\n", err)
	} else {
		for i, val := range maps {
			fmt.Printf("  Map %d is: %+v\n", i, val)
		}
	}

	fmt.Println("\nReading sub array")
	subs, err := c.Childs("maps")
	if err != nil {
		fmt.Printf("  Error getting sub array: %+v\n", err)
	} else {
		for i, sub := range subs {
			foo, _ := sub.String("foo")
			fmt.Printf("  Sub %d is: %+v and value of foo is \"%s\"\n", i, sub.Root(), foo)
		}
	}
}
