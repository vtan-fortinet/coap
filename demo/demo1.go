package main

import (
    "fmt"
    "coap"
)


type myArg struct {
    Uname string    `-nNAME --name
                    !what is your name`
    Passwd string   `-pPASS --passwd
                    user password`
    Act string      `---ACT
                    ! user action and arg
                    -u --upload
                    upload file with username/password
                    -d --download
                    download and save into file`
}


func main() {
    m := myArg{Passwd: "123456"}
    args := coap.Parse(&m)
    fmt.Printf("-n = %s, -p = '%s'\n", m.Uname, m.Passwd)
    fmt.Printf("act = %s\n", m.Act)
    fmt.Printf("args = %v\n", args)
}
