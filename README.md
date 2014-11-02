[![Build Status](https://travis-ci.org/ukautz/mappath.svg?branch=master)](https://travis-ci.org/ukautz/mappath)

MapPath
=======

[Go](http://golang.org/) library for convenient read-access to data structures.

Purpose and Scope
-----------------

This is not an [XPATH](http://en.wikipedia.org/wiki/XPath) implementation for Go, but a simple path interface
to structured data files, like JSON or YAML. Think accessing configuration files or the like.

Documentation
-------------

GoDoc can be [found here](http://godoc.org/github.com/ukautz/mappath)

### Installation

```bash
$ go get github.com/ukautz/mappath
```

### Usage

This package needs at least Go 1.1. Import package with

```go
import "github.com/ukautz/mappath"
```

Then you can do

```go
mp, err := mappath.FromJsonFile("path/to.json")
if err != nil {
    panic(err)
}

str, err := mp.GetString("foo/bar/0/hello")
```

The above assumes a simple JSON file:

```json
{
    "foo": {
        "bar": [
            {
                "hello": "world"
            }
        ]
    }
}
```

### Intialization

The library works on any `map[string]interface{}` structures. Those can come from JSON/YAML files or anywhere else. The above example illustrates how to get easily started with an existing JSON file. Here is how you do it Go-only:

```go
mp := mappath.NewMapPath(map[string]interface{}{
    "foo": map[string]interface{}{
        "bar": map[string]interface{}{
            "baz": "hello",
        },
    },
})
v, _ := mp.GetString("foo/bar/baz")
fmt.Printf("Say %s world\n", v)
```

### Accessing data

```go
mp := mappath.NewMapPath(source)

// get an interface{}
result, err := mp.Get("the/path")

// get an int value
// assumes a structure like: {"the":{"path":123}}
result, err = mp.GetInt("the/path")

// get an array of int values
// assumes a structure like: {"the":{"path":[123, 234]}}
result, err = mp.GetInts("the/path")

// get a float value
// assumes a structure like: {"the":{"path":123.1}}
result, err = mp.GetFloat("the/path")

// get an array float values
// assumes a structure like: {"the":{"path":[123.1, 234.9]}}
result, err = mp.GetFloats("the/path")

// get a string value
// assumes a structure like: {"the":{"path":"foo"}}
result, err = mp.GetString("the/path")

// get an array of string values
// assumes a structure like: {"the":{"path":["foo","bar"}}
result, err = mp.GetString("the/path")

// get a map value
// assumes a structure like: {"the":{"path":{"foo":"bar"}}
result, err = mp.GetMap("the/path")

// get an array of map values
// assumes a structure like: {"the":{"path":[{"foo":"bar1"},{"foo":"bar2"}]}
result, err = mp.GetMap("the/path")
```

### Using sub structures

For example, when iterating above an a structure like the following

```json
{
    "users":[
        {
            "name":"Mr Kirk"
        },
        {
            "name":"Mr Spock"
        },
        {
            "name":"Jean-Luc"
        }
    ]
}
```
Here is how:

```go
subs, err := mp.GetSubs("users")
for _, user := range subs {
    fmt.Printf("Say hello to %s\n", user.GetString("name"))
}
```

### Error handling

**`mappath.NotFoundError`**

When accessing a path with any `Get` method, the result contains the (appropriate) `nil` value
if the path cannot be found.

**`mappath.InvalidTypeError`**

If you use the type specific getter you could also try getting a value which cannot be convert. For example
when trying to get an `int` of a string value like `foo bar`, or when you try to get a `string` but the
path is actually a sub-structure.

**`mappath.UnsupportedTypeError`**

Used when you try to get an array of a not supported type. At the moment, those are `int`, `float64`, `string` and `map[string]interface{}`.

### Convenience: Fallback values

Since I developed this library mainly for working with complex configuration files it's a common use-case to provide
a _fallback_ value, which is used if nothing is found at the given path. If you provide a fallback value, than no
`NotFoundError` will be returned.

```go
// returns "Some Fallback" if the path does not exist
result, err := mp.GetString("the/path", "Some Fallback")
```

## See also

* A completely different approach by Yasuyuki YAMADA to simplify access to JSON files: https://github.com/yasuyuky/jsonpath
