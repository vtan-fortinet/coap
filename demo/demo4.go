package main

import (
    //"os"
    "fmt"
    //"time"
    "github.com/pastebt/coap"
)


type GRP struct {
    Sel string
    Val bool
}


type myArg struct {
    Act *GRP    `---
                 different action
                 -s --start
                 start
                 -e --end
                 end`
    B bool      `-b --bool
                test bool`
}


func main() {
    m := myArg{}
    args := coap.Parse(&m)
    fmt.Printf("Act = %v\n", m.Act)
    fmt.Printf("args = %v\n", args)
}
