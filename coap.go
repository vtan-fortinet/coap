package coap

import (
    "os"
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
}


func (c *COAP)Parse() { c.ParseA(os.Args[1:]) }
func (c *COAP)ParseA(args []string) {
    c.init()
}


func (c *COAP)Help(err_msg string) {
    c.init()
}
