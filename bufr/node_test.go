package bufr

import (
    "testing"
    "fmt"
)

func TestBlock_Add(t *testing.T) {
    var b Block
    //var b1 Block


    fmt.Println(b.Members())
    fmt.Printf("%+v\n", b)
}