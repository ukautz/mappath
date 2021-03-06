[![Build Status](https://travis-ci.org/ukautz/mappath.svg?branch=master)](https://travis-ci.org/ukautz/mappath)

MapPath
=======

[Go](http://golang.org/) library for convenient read-access to data structures.

Purpose and Scope
-----------------

This is not an [XPATH](http://en.wikipedia.org/wiki/XPath) implementation for Go, but a simple path interface
to structured data files, like JSON or YAML. Think accessing configuration files or the like.

API Changes
-----------

### v1 -> v2

* Removed the "Get" prefix of all methods, so former `mappath.GetInt("foo")` becomes `mappath.Int("foo")`. The outlier is `GetSub` which is now `Child`.
* Added `V`alue-getter with scalar response, eg `mappath.IntV("foo")` has the return signatur of `int`, while `mappath.Int("foo")` still has `(int, error)`. The `V`-getter return the `nil` value, on error

Documentation
-------------

GoDoc can be [found here](http://godoc.org/github.com/ukautz/mappath)

### Installation

```bash
$ go get gopkg.in/ukautz/mappath.v2
```

### Usage

This package needs at least Go 1.1. Import package with

```go
import "gopkg.in/ukautz/mappath.v2"
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

The library works on `map[string]interface{}` structures. Those can come from JSON/YAML files or anywhere else. The above example illustrates how to get easily started with an existing JSON file. Here is how you do it Go-only:

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
result, err = mp.GetStrings("the/path")

// get a map value
// assumes a structure like: {"the":{"path":{"foo":"bar"}}
result, err = mp.GetMap("the/path")

// get an array of map values
// assumes a structure like: {"the":{"path":[{"foo":"bar1"},{"foo":"bar2"}]}
result, err = mp.GetMaps("the/path")
```

### Using sub structures

For example, when iterating a structure like the following

```json
{
    "users":[
        {
            "name":"Cpt Kirk"
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
Here is how you can conveniently access the nested maps:

```go
subs, err := mp.GetSubs("users")
for _, user := range subs {
    fmt.Printf("Say hello to %s\n", user.GetString("name"))
}
```

### Error handling

**`mappath.NotFoundError`**

Returned if the accessed path does not exist. The result will contain the appropriate `nil` value.

**`mappath.InvalidTypeError`**

Returned if you get a path which exists but contains a value which can neither be converted nor parsed. For example:
trying to get the `int` value of the string `foo bar`.

**`mappath.UnsupportedTypeError`**

Returned on array getter. The currently supported types are: `int`, `float64`, `string` and `map[string]interface{}`.

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
