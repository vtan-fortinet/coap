package main

import (
    "fmt"
    "../../coap"
)


type myArg struct {
    Uname string    `-nNAME --name
                    !what is your name`
    Passwd string   `-pPASS --passwd
                    user password`
    Lvl int         `-lLEVEL --level 2|LVL
                    [1, 2, 3]
                    comrpess level`
    Act string      `---ACT
                    ! user action and arg
                    -u --upload
                    upload file with username/password
                    -d --download
                    download and save into file`
    Bool bool       `-b 
                    test bool`
}


func main() {
    m := myArg{Passwd: "123456"} //, Lvl: 2}
    args := coap.Parse(&m)
    fmt.Printf("-n = %s, -p = '%s'\n", m.Uname, m.Passwd)
    fmt.Printf("-l = %d, act = %s\n", m.Lvl, m.Act)
    fmt.Printf("args = %v\n", args)
}
