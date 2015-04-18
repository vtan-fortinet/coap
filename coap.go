package coap

import (
    "os"
    "fmt"
    "reflect"
)


type oaItem struct {    // option, argument item
    short   string
    long    string
    must    bool
    help    string
}


type GRP struct {       // group
    sel     string
    val     string
}


type COAP struct {
    items   []oaItem
}


func (c *COAP)init() {
    if c.items != nil { return }
    //panic("init")
    c.items = make([]oaItem, 0, 10)
    st := reflect.TypeOf(*c)
    fmt.Println(st)
    fmt.Println(st.Field(0))
    fmt.Println(st.Field(1))
    fmt.Println(st.Field(2))
    fmt.Println(st.NumField())
}


func (c *COAP)Parse() { c.ParseArgs(os.Args[1:]) }
func (c *COAP)ParseArgs(args []string) {
    c.init()
}


func (c *COAP)Help() { c.HelpMsg("") }
func (c *COAP)HelpMsg(msg string) {
    c.init()
}


func Parse(arg interface{}) {
    fmt.Println("Parse(arg interface{})", arg)
    st := reflect.TypeOf(arg)
    fmt.Println(st)
    fmt.Println(st.NumField())
}
