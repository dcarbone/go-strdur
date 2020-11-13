# go-strdur

A highly inefficient implementation of `time.Duration` based on a string, rather than an int64.  It is designed purely
soo that hcl configuration files can contain a string representation of a duration for config values, rather than an
int64.

## Example

```go
package main

import (
    "fmt"
    
    "github.com/dcarbone/go-strdur"
    "github.com/hashicorp/hcl/v2/hclsimple"
)

const exampleConfig = `
duration_value = "24h"
`

type MyConfig struct {
    DurationValue strdur.StringDuration `hcl:"duration_value"`
}

func main() {
    myCnf := new(MyConfig)
    if err := hclsimple.Decode("example.hcl", []byte(exampleConfig), nil, new(MyConfig)); err != nil {
        panic(fmt.Sprintf("Error decoding hcl: %v", err))
    }
    fmt.Println(myCnf.DurationValue.String())
}
```

## Explanation
I created this type specifically because, as of the time of this writing, https://github.com/hashicorp/hcl does not have
a great way to handle embedded types.

Given this example:

*GO*:
```go
package main

import (
    "fmt"
    "time"

    "github.com/hashicorp/hcl/v2/hclsimple"
)

const confValue = `
duration_value = "24h"
`

type MyDuration time.Duration

type MyConfig struct {
    DurationValue MyDuration `hcl:"duration_value"`
}

func main() {
    if err := hclsimple.Decode("example.hcl", []byte(confValue), nil, new(MyConfig)); err != nil {
        panic(fmt.Sprintf("Error decoding hcl: %v", err))
    }
}
```

The above will always fail.  Its possible I am missing something in the
[cty](https://pkg.go.dev/github.com/zclconf/go-cty/cty) package which can handle this, but I haven't been able to find
it.

I offer zero guarantees of performance on this type as it is intended entirely to be used as a value in a config object.