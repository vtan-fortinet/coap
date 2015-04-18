package main

import (
    "os"
    "fmt"
    "coap"
)


type myArg struct {
    //coap.COAP
    name string     `-n --name
                    what is your name`
}


func main() {
    fmt.Println("demo1")
    m := myArg{}
    //m.Parse()
    //m.Help()
    coap.Help(&m)
    //coap.Parse(&m)
    os.Exit(1)
}
