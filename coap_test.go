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
    if ! oa.HasDft || oa.MsgDft != "" || oa.StrDft != "10" {
        tst.Error("failed to parse default5", oa)
    }

    i = 10
    oa = oaItem{}
    oa.initDefault(reflect.ValueOf(i), "|TEN")
    if ! oa.HasDft || oa.MsgDft != "TEN" || oa.StrDft != "10" {
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


func TestInitGrp1(tst *testing.T) {
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


type B struct {
    g0 *G    `---GRP 
              -o --open
              group open help
            `

    g1 *G    `---GRP b
              -a --abc
              group abc help
              -b --bac
              group bac help
            `
}


func TestInitGrp2(tst *testing.T) {
    v := reflect.ValueOf(B{})
    t := v.Type()
    oa := oaItem{}
    oa.init(t.Field(0), v.Field(0))
    if oa.Grp == nil || len(oa.Grp) != 1 || oa.Vname != "GRP" || oa.Must {
        tst.Error("failed to init grp0", oa)
    }
    if oa.HasDft || len(oa.HelpLs) > 0 || oa.Must {
        tst.Error("failed to init grP0", oa.HasDft, oa.HelpLs)
    }
    if oa.Grp[0].Short != "o" || oa.Grp[0].Long != "open" {
        tst.Error("failed to init GRP0", oa.Grp[0])
    }
    if len(oa.Grp[0].HelpLs) < 1 || oa.Grp[0].HelpLs[0] != "group open help" {
        tst.Error("failed to init GRH0", oa.Grp[0])
    }

    oa = oaItem{}
    oa.init(t.Field(1), v.Field(1))
    if oa.Grp == nil || len(oa.Grp) != 2 || oa.Vname != "GRP" || oa.Must {
        tst.Error("failed to init grp1", oa)
    }
    if ! oa.HasDft || len(oa.HelpLs) > 0 || oa.Must {
        tst.Error("failed to init grP1", oa.HasDft, oa.HelpLs)
    }
    if oa.Grp[0].Short != "a" || oa.Grp[0].Long != "abc" {
        tst.Error("failed to init GRP1", oa.Grp[0])
    }
    if len(oa.Grp[0].HelpLs) < 1 || oa.Grp[0].HelpLs[0] != "group abc help" {
        tst.Error("failed to init GRH1", oa.Grp[0])
    }
    if oa.Grp[1].Short != "b" || oa.Grp[1].Long != "bac" {
        tst.Error("failed to init GRP1", oa.Grp[1])
    }
    if len(oa.Grp[1].HelpLs) < 1 || oa.Grp[1].HelpLs[0] != "group bac help" {
        tst.Error("failed to init GRH1", oa.Grp[1])
    }
}


func TestGet_next(tst *testing.T) {
    args := []string{"-a", "b"}
    o, a := get_next(0, args)
    if *o != "-a" || *a != "b" {
        tst.Error("failed to get_next1", o, a)
    }

    args = []string{"-a"}
    o, a = get_next(0, args)
    if *o != "-a" || a != nil {
        tst.Error("failed to get_next2", o, a)
    }

    args = []string{"-a", "-b"}
    o, a = get_next(0, args)
    if *o != "-a" || a != nil {
        tst.Error("failed to get_next3", o, a)
    }

    args = []string{"-a", "-b", "c"}
    o, a = get_next(1, args)
    if *o != "-b" || *a != "c" {
        tst.Error("failed to get_next4", o, a)
    }
}


func TestSetValuei(tst *testing.T) {
    var i int
    var i8 int8
    var i16 int16
    var i32 int32
    var i64 int64

    v := reflect.ValueOf(&i).Elem()
    got, err := setValue(&v, "-12")
    if err != "" || i != -12 || got != 2 {
        tst.Error("failed to setValuei")
    }

    v = reflect.ValueOf(&i8).Elem()
    got, err = setValue(&v, "-123")
    if err != "" || i8 != -123 || got != 2 {
        tst.Error("failed to setValuei8")
    }

    v = reflect.ValueOf(&i16).Elem()
    got, err = setValue(&v, "-12345")
    if err != "" || i16 != -12345 || got != 2 {
        tst.Error("failed to setValuei16")
    }

    v = reflect.ValueOf(&i32).Elem()
    got, err = setValue(&v, "-1234578")
    if got != 2 || err != "" || i32 != -1234578 {
        tst.Error("failed to setValuei32")
    }

    v = reflect.ValueOf(&i64).Elem()
    got, err = setValue(&v, "-123457890")
    if got != 2 || err != "" || i64 != -123457890 {
        tst.Error("failed to setValuei64")
    }
}


func TestSetValueu(tst *testing.T) {
    var u uint
    v := reflect.ValueOf(&u).Elem()
    got, err := setValue(&v, "12")
    if got != 2 || err != "" || u != 12 {
        tst.Error("failed to setValueu")
    }

    var u8 uint8
    v = reflect.ValueOf(&u8).Elem()
    got, err = setValue(&v, "123")
    if got != 2 || err != "" || u8 != 123 {
        tst.Error("failed to setValueu8")
    }

    var u16 uint16
    v = reflect.ValueOf(&u16).Elem()
    got, err = setValue(&v, "12345")
    if got != 2 || err != "" || u16 != 12345 {
        tst.Error("failed to setValueu16")
    }

    var u32 uint32
    v = reflect.ValueOf(&u32).Elem()
    got, err = setValue(&v, "1234578")
    if got != 2 || err != "" || u32 != 1234578 {
        tst.Error("failed to setValueu32")
    }

    var u64 uint64
    v = reflect.ValueOf(&u64).Elem()
    got, err = setValue(&v, "123457890")
    if got != 2 || err != "" || u64 != 123457890 {
        tst.Error("failed to setValueu64")
    }
}


func TestSetValuefs(tst *testing.T) {
    var f32 float32
    v := reflect.ValueOf(&f32).Elem()
    got, err := setValue(&v, "-12.34")
    if got != 2 || err != "" || f32 != -12.34 {
        tst.Error("failed to setValuef32")
    }

    var f64 float64
    v = reflect.ValueOf(&f64).Elem()
    got, err = setValue(&v, "123456.789")
    if got != 2 || err != "" || f64 != 123456.789 {
        tst.Error("failed to setValuef64")
    }

    var s string
    v = reflect.ValueOf(&s).Elem()
    got, err = setValue(&v, "abc123456.789")
    if got != 2 || err != "" || s != "abc123456.789" {
        tst.Error("failed to setValues")
    }

    var b bool
    v = reflect.ValueOf(&b).Elem()
    got, err = setValue(&v, "true")
    if got != 1 || err != "" || ! b {
        tst.Error("failed to setValueb", got, b)
    }
}


func rp(tst *testing.T, msg string) {
    r := recover()
    if msg != r {
        tst.Errorf("failed to get panic [%s], got [%s]\n", msg, r)
    }
}
func TestRp(tst *testing.T) {
    defer rp(tst, "I am panic")
    panic("I am panic")
}


func rp2(msg *string) {
    *msg = recover().(string)
}
func TestRp2(tst *testing.T) {
    var msg string
    func() {
    //defer func(m *string) { *m = recover().(string) }(&msg)
    defer rp2(&msg)
    panic("I am panic")
    }()
    if msg != "I am panic" {
        tst.Error("failed to get panic", msg)
    }
}

