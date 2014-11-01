# MapPath

Library for convenient access to data structures.

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

Assuming a simple JSON file:

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

## Extended example

TODO

