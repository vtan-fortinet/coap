package main

import (
    //"os"
    "fmt"
    "time"
    "../../coap"
)


type myArg struct {
    Date string    `-dDate --date
                    !report date`
    Tst bool        `--test
                    test status`
}


func main() {
    m := myArg{Date: time.Now().AddDate(0, 0, -1).Format("2006-01-02")}
    //m := myArg{}
    coap.RegValFunc(&m, "d", func (i interface{}) string {
        m := i.(*myArg)
        _, err := time.Parse("2006-01-02", m.Date)
        if err != nil {
            return "Wrong date format, should like '2006-01-02'"
        }
        return ""
    })
    args := coap.Parse(&m)
    fmt.Printf("-d = %s\n", m.Date)
    fmt.Printf("args = %v\n", args)
}
