package main

import (
//    "os"
    "fmt"
    "coap"
)


type myArg struct {
    Uname string    `-nNAME --name
                    !what is your name`
    Passwd string   `-pPASS --passwd
                    user password`
}


func main() {
    m := myArg{Passwd: "123456"}
    args := coap.Parse(&m)
    //fmt.Println("uname =", m.Uname)
    //fmt.Println(os.Args)
    //fmt.Fprint(os.Stderr, "err msg\n")
    fmt.Printf("-n = %s, -p = '%s', args = %v\n", m.Uname, m.Passwd, args)
    //os.Exit(1)
}
