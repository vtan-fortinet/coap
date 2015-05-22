package main

import (
    "os"
    "fmt"
    "time"
    "github.com/pastebt/coap"
)


type myArg struct {
    Date string    `-dDate --date
                    !report date`
}


func main() {
    //m := myArg{Date: time.Now().AddDate(0, 0, -1).Format("2006-01-02")}
    m := myArg{}
    args := coap.Parse(&m)
    _, err := time.Parse("2006-01-02", m.Date)
    if err != nil {
        coap.HelpMsg(&m, "Wrong date format, should like '2006-01-02'", os.Stderr)
        os.Exit(1)
    }
    fmt.Printf("-d = %s\n", m.Date)
    fmt.Printf("args = %v\n", args)
}
