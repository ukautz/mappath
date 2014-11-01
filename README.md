[![Build Status](https://travis-ci.org/ukautz/mappath.svg?branch=master)](https://travis-ci.org/ukautz/mappath)

# MapPath

[Go](http://golang.org/) library for convenient read-access to data structures.

## Purpose & Scope

This is not an [XPATH](http://en.wikipedia.org/wiki/XPath) implementation for Go, but a simple path interface
 to structured data files, like JSON or YAML. Think accessing configuration files or the like.

## Installation

```bash
$ go get github.com/ukautz/mappath
```

## Usage

This package needs at least Go 1.1. Import package with

```go
import "github.com/ukautz/mappath
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

### Access methods

```go
mp := mappath.NewMapPath(source)

# get an interface{}
result, err := mp.Get("the/path")

# get a string
result, err = mp.GetString("the/path")

# get an int
result, err = mp.GetInt("the/path")

# get a float
result, err = mp.GetFloat("the/path")

# get a map[string]interface{} (i.e. a sub-structre)
result, err = mp.GetMap("the/path")
```

### Error handling

**`mappath.NotFoundError`**

When accessing a path with any `Get` method, the result contains the (appropriate) `nil` value
if the path cannot be found.

**`mappath.InvalidTypeError`**

If you use the type specific getter you could also try getting a value which cannot be convert. For example
when trying to get an `int` of a string value like `foo bar`, or when you try to get a `string` but the
path is actually a sub-structure.

### Convenience: Fallback values

If you work with config files, you usually have some kind of default you want to use if a config value does not exist.

```go
result, err := mp.GetString("the/path", "Some Fallback")
```

## See also

* A completely different approach by Yasuyuki YAMADA to simplify access to JSON files: https://github.com/yasuyuky/jsonpath
