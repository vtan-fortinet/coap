package coap

import (
    "testing"
    "reflect"
)


func TestInitHelp1(tst *testing.T) {
    oa := oaItem{}
    oa.initHelp(`!this is a help msg`)
    if ! oa.Must {
        tst.Error("should be must")
    }
    if len(oa.HelpLs) != 1 || oa.HelpLs[0] != "this is a help msg" {
        tst.Error("help error parse first line")
    }
}


func TestInitHelp2(tst *testing.T) {
    oa := oaItem{}
    oa.initHelp(`!!this is a help msg`)
    if oa.Must {
        tst.Error("should not be must")
    }
    if len(oa.HelpLs) != 1 || oa.HelpLs[0] != "!this is a help msg" {
        tst.Error("help error parse first line")
    }
}


func TestInitHelp3(tst *testing.T) {
    oa := oaItem{}
    oa.initHelp(`! !this is a help msg`)
    if ! oa.Must {
        tst.Error("should be must")
    }
    if len(oa.HelpLs) != 1 || oa.HelpLs[0] != "!this is a help msg" {
        tst.Error("help error parse first line")
    }
}


func TestInitHelp4(tst *testing.T) {
    oa := oaItem{}
    oa.initHelp(`this is a help msg`)
    oa.initHelp(`Here is second line`)
    if oa.Must {
        tst.Error("should not be must")
    }
    //if oa.HelpLs != []string{"this is a help msg", "Here is second line"} {
    if len(oa.HelpLs) != 2 || oa.HelpLs[0] != "this is a help msg" || oa.HelpLs[1] != "Here is second line" {
        tst.Error("help error parse lines", oa.HelpLs)
    }
}


func TestInitOpts(tst *testing.T) {
    oa := oaItem{}
    oa.initOpts("-n")
    if oa.Short !=  "n" { tst.Error("failed parse short") }

    oa = oaItem{}
    oa.initOpts("--name")
    if oa.Short !=  "" || oa.Long != "name" { tst.Error("failed parse long") }

    oa = oaItem{}
    oa.initOpts("-n --name")
    if oa.Short !=  "n" || oa.Long != "name" { tst.Error("failed parse short and long") }

    oa = oaItem{}
    oa.initOpts("-nNAME")
    if oa.Short !=  "n" || oa.Long != "" || oa.Vname != "NAME"{
        tst.Error("failed parse short and long")
    }

    oa = oaItem{}
    oa.initOpts("-nNAME    --name")
    if oa.Short !=  "n" || oa.Long != "name" || oa.Vname != "NAME"{
        tst.Error("failed parse short and long")
    }

    oa = oaItem{}
    d := oa.initOpts("-nNAME    --name   dft|DEFATLT")
    if oa.Short !=  "n" || oa.Long != "name" || oa.Vname != "NAME" || d != "dft|DEFATLT" {
        tst.Error("failed parse short and long")
    }
}


func TestInitDefault(tst *testing.T) {
    var i int
    oa := oaItem{}
    oa.initDefault(reflect.ValueOf(i), "10|Ten")
    if ! oa.HasDft || oa.MsgDft != "Ten" || oa.StrDft != "10" {
        tst.Error("failed to parse default1", oa)
    }

    oa = oaItem{}
    oa.initDefault(reflect.ValueOf(i), "10")
    if ! oa.HasDft || oa.MsgDft != "" || oa.StrDft != "10" {
        tst.Error("failed to parse default2", oa)
    }

    oa = oaItem{}
    oa.initDefault(reflect.ValueOf(i), "|Ten")
    if oa.HasDft || oa.MsgDft != "Ten" || oa.StrDft != "" {
        tst.Error("failed to parse default3", oa)
    }

    oa = oaItem{}
    oa.initDefault(reflect.ValueOf(i), "")
    if oa.HasDft || oa.MsgDft != "" || oa.StrDft != "" {
        tst.Error("failed to parse default4", oa)
    }

    i = 10
    oa = oaItem{}
    oa.initDefault(reflect.ValueOf(i), "")
    if ! oa.HasDft || oa.MsgDft != "" || oa.StrDft != "" {
        tst.Error("failed to parse default5", oa)
    }

    i = 10
    oa = oaItem{}
    oa.initDefault(reflect.ValueOf(i), "|TEN")
    if ! oa.HasDft || oa.MsgDft != "TEN" || oa.StrDft != "" {
        tst.Error("failed to parse default6", oa)
    }
}


type G struct {
    sel string
    val bool
}


type A struct {
    g0 *G    `---GRP `
    g1 *G    `---GRP b
            `
    g2 *G    `---GRP b|GROUP
            `
    g3 *G    `---GRP b|GROUP
              group help
            `
    g4 *G    `---GRP b|GROUP
              !group help must `
}


func TestInitGrp(tst *testing.T) {
    v := reflect.ValueOf(A{})
    t := v.Type()
    oa := oaItem{}
    oa.init(t.Field(0), v.Field(0))
    if oa.Grp == nil || len(oa.Grp) > 0 || oa.Vname != "GRP" || oa.Must {
        tst.Error("failed to init grp0", oa)
    }
    if oa.HasDft || len(oa.HelpLs) > 0 || oa.Must {
        tst.Error("failed to init grP0", oa.HasDft, oa.HelpLs)
    }

    oa = oaItem{}
    oa.init(t.Field(1), v.Field(1))
    if oa.Grp == nil || len(oa.Grp) > 0 || oa.Vname != "GRP" || oa.Must {
        tst.Error("failed to init grp1", oa)
    }
    if ! oa.HasDft || oa.StrDft != "b" || oa.MsgDft != "" || len(oa.HelpLs) != 1 {
        tst.Error("failed to init grP1", oa.HasDft, oa.StrDft, oa.MsgDft, oa.HelpLs)
    }

    oa = oaItem{}
    oa.init(t.Field(2), v.Field(2))
    if oa.Grp == nil || len(oa.Grp) > 0 || oa.Vname != "GRP" || oa.Must {
        tst.Error("failed to init grp2", oa)
    }
    if ! oa.HasDft || oa.StrDft != "b" || oa.MsgDft != "GROUP"{
        tst.Error("failed to init grP2", oa)
    }

    oa = oaItem{}
    oa.init(t.Field(3), v.Field(3))
    if oa.Grp == nil || len(oa.Grp) > 0 || oa.Vname != "GRP" || oa.Must {
        tst.Error("failed to init grp3", oa)
    }
    if ! oa.HasDft || oa.StrDft != "b" || oa.MsgDft != "GROUP" {
        tst.Error("failed to init grP3", oa.HasDft, oa.StrDft, oa.MsgDft)
    }
    if len(oa.HelpLs) < 1 || oa.HelpLs[0] != "group help"{
        tst.Error("failed to init grH3", oa.HelpLs)
    }

    oa = oaItem{}
    oa.init(t.Field(4), v.Field(4))
    if oa.Grp == nil || len(oa.Grp) > 0 || oa.Vname != "GRP" || ! oa.Must {
        tst.Error("failed to init grp3", oa)
    }
    if ! oa.HasDft || oa.StrDft != "b" || oa.MsgDft != "GROUP" {
        tst.Error("failed to init grP3", oa.HasDft, oa.StrDft, oa.MsgDft)
    }
    if len(oa.HelpLs) < 1 || oa.HelpLs[0] != "group help must"{
        tst.Error("failed to init grH3", oa.HelpLs)
    }

}
