package main

import (
    "fmt"
    "coap"
)


type myArg struct {
    Act string   `---
                 !different action
                 -s --start
                 start
                 -e --end
                 end`
    Vs []bool    `-v
                 test bool`
}


func main() {
    m := myArg{}
    args := coap.Parse(&m)
    fmt.Printf("Act = %v\n", m.Act)
    fmt.Printf("Vs = %v\n", m.Vs)
    fmt.Printf("args = %v\n", args)
}
