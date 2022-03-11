# golang-edit-prompt
[![Go Reference](https://pkg.go.dev/badge/github.com/ibrt/golang-edit-prompt.svg)](https://pkg.go.dev/github.com/ibrt/golang-edit-prompt)
![CI](https://github.com/ibrt/golang-edit-prompt/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/ibrt/golang-edit-prompt/branch/main/graph/badge.svg?token=BQVP881F9Z)](https://codecov.io/gh/ibrt/golang-edit-prompt)

Run visudo-like text editing prompts from a Go CLI program.

### Basic Example

```go
package main

import (
    "github.com/ibrt/golang-edit-prompt/editz"
)

func main() {
    err := editz.Edit("hello.txt", func(contents []byte) error {
        // inspect the changed contents, return nil if OK, error otherwise
        return nil
    })
    if err != nil {
        panic(err)
    }
}
```

### Developers

Contributions are welcome, please check in on proposed implementation before sending a PR. You can validate your changes using the `./test.sh` script.
