package coap

import (
    "os"
)


type oAI struct {       // option, argument item
    short   string
    long    string
    must    bool
    help    string
}


type gRP struct {       // group
    sel     string
    val     string
}


type COAP struct {
    items   []oAI
}


func (c *COAP)Parse() { c.ParseA(os.Args[1:]) }
func (c *COAP)ParseA(args []string) {

}


func (c *COAP)Help(err_msg string) {

}
