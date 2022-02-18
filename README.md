## go-delta-sib - A Go package and Custom library for generating delta using snappy compression
[![godoc](https://godoc.org/github.com/DTSL/go-delta-sib?status.svg)](https://godoc.org/github.com/DTSL/go-delta-sib)

## Suggestions:

- Works best on text files, database dumps and any other files with lots of
  repeating patterns and few changes between updates.

- Don't compress bytes returned by Delta.Bytes() because they are already
  compressed using Snappy compression.

## Demonstration:

```go
package main

import (
    "fmt"
    "github.com/DTSL/go-delta-sib"
)

func main() {
    fmt.Print("Binary delta update demo:\n\n")

    // The original data (20 bytes):
    var source = []byte("quick brown fox, lazy dog, and five boxing wizards")
    fmt.Print("The original is:", "\n", string(source), "\n\n")

    // The updated data containing the original and new content (82 bytes):
    var target = []byte(
        "The quick brown fox jumps over the lazy dog. " +
        "The five boxing wizards jump quickly.",
    )
    fmt.Print("The update is:", "\n", string(target), "\n\n")

    var dbytes []byte
    {
    	// Use Make() to generate a compressed patch from source and target
    	var d = delta.Make(source, target)
    	
    	// Convert the delta to a slice of bytes (e.g. for writing to a file)
    	dbytes = d.Bytes()
    }

    // Create a Delta from the byte slice
    var d = delta.Load(dbytes)

    // Apply the patch to source to get the target
    // The size of the patch is much shorter than target.
    var target2, err = d.Apply(source)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Print("Patched:", "\n", string(target2), "\n\n")
} //                                                                        main
```
