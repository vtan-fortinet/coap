package main

import (
    "fmt"
    "coap"
)


type myArg struct {
    coap.COAP
    name string     `-n --name
                    what is your name`
}


func main() {
    fmt.Println("demo1")
}
