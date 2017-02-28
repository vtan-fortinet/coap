package main

import (
    "coap"
)


type myArg struct {
    Date string `-dDATE
                date for save`

    G1 string   `---GRP
                !Help for this group
                -b --begin
                help for begin: begin the service
                -e --end
                end the service`
}


func main() {
    mo := &myArg{}
    a  := coap.Parse(mo)
    //println("sel=", mo.InDn.Sel, ", val=", mo.InDn.Val)
    println("G1=", mo.G1)
    println("A=", a)
}
