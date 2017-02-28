package main

import (
    "fmt"
    "github.com/pastebt/coap"
)


type myArg struct {
    Act string   `---
                 !different action
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
