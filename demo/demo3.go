package main

import (
    "coap"
)

type Grp struct {
    Sel string  // first must be string
    Val string  // type is whatever can be converted from string
}


type myArg struct {
    //InDn *Grp   `---DIRNAME
    //            Input/Verify directory
    //            -i --input
    //            input directory
    //            -v --verify
    //            verify directory`

    Date string `-dDATE
                date for save`

    G1 string   `---GRP
                !Help for this group
                -b --begin
                help for begin: begin the service
                -e --end
                end the service`
    Host string `-h --host
                hostname`
}


func main() {
    mo := &myArg{}
    a  := coap.Parse(mo)
    //println("sel=", mo.InDn.Sel, ", val=", mo.InDn.Val)
    println("G1=", mo.G1)
    println("A=", a)
}
