package main

import (
    "os"
    "fmt"
    "coap"
)


type myArg struct {
    //coap.COAP
    Uname string    `-nNAME --name
                    !what is your name`
}


func main() {
    //fmt.Println("demo1")
    m := myArg{}
    //m.Parse()
    //m.Help()
    //coap.Help(&m)
    coap.Parse(&m)
    //fmt.Println("uname =", m.Uname)
    //fmt.Println(os.Args)
    fmt.Printf("-n = %s\n", m.Uname)
    os.Exit(0)
}
