package main

import (
    "os"
    "fmt"
    "coap"
)


type myArg struct {
    //coap.COAP
    Uname string    `-n --name
                    what is your name`
}


func main() {
    fmt.Println("demo1")
    m := myArg{}
    //m.Parse()
    //m.Help()
    coap.Help(&m)
    //coap.Parse(&m)
    fmt.Println("uname =", m.Uname)
    os.Exit(0)
}
