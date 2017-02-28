package coap

import (
    "fmt"
    "testing"
    "reflect"
    "strconv"
)


func reset() { infos = make(map[uintptr]*oaInfo) }


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


func TestDft1(tst *testing.T) {
    type dft struct {
        S string `-s --string ""|empty
                  !this is em`
    }
    d := &dft{}
    msg, ps := ParseArg(d, []string{})
    if msg != "Missed option -s" || len(ps) != 0 {
        tst.Error(msg)
        tst.Error(ps)
    }

    d = &dft{}
    msg, ps = ParseArg(d, []string{"-s"})
    if msg != "" || len(ps) != 0 || d.S != "" {
        tst.Error(msg)
        tst.Error(ps)
    }

    d = &dft{}
    msg, ps = ParseArg(d, []string{"-s", "str"})
    if msg != "" || len(ps) != 0 || d.S != "str" {
        tst.Error(msg)
        tst.Error(ps)
    }
}


func TestInitGrp3(tst *testing.T) {
    type g31 struct {
        G3 string `---GRP3
                   !help for g3
                   -a --aa
                   help for g3a
                   -b --bb
                   help for g3b`
    }

    g := &g31{}
    msg, ps := ParseArg(g, make([]string, 0))
    if msg != "Missed option -a|-b" || len(ps) != 0 {
        tst.Error(msg)
        tst.Error(ps)
    }

    g = &g31{}
    msg, ps = ParseArg(g, []string{"-a"}[:])
    if msg != "option -a need parameter" || len(ps) != 0 {
        tst.Error(msg)
        tst.Error(ps)
    }

    g = &g31{}
    msg, ps = ParseArg(g, []string{"-a", "aa"}[:])
    if msg != "" || len(ps) != 0 {
        tst.Error(msg)
        tst.Error(ps)
    }

    g = &g31{}
    msg, ps = ParseArg(g, []string{"-a", "aa", "cc"}[:])
    if msg != "" || len(ps) != 1 || ps[0] != "cc" {
        tst.Error(msg)
        tst.Error(ps)
    }
}


func TestInitGrp4(tst *testing.T) {
    type g41 struct {
        G4 string `---GRP4
                   !help for g4
                   -a 
                   help for g3a
                   -b
                   help for g3b`
        S string `-s --string
                  string i`
    }

    g := &g41{}
    msg, ps := ParseArg(g, []string{})
    if msg != "Missed option -a|-b" || len(ps) != 0 {
        tst.Error(msg)
        tst.Error(ps)
    }

    g = &g41{}
    msg, ps = ParseArg(g, []string{"-a", "aa"})
    if msg != "" || len(ps) != 0 || g.G4 != "a aa" {
        tst.Error(msg)
        tst.Error(ps)
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
    if *o != "-a" || *a != "-b" {
        tst.Error("failed to get_next3", o, a)
    }

    args = []string{"-a", "-b", "c"}
    o, a = get_next(1, args)
    if *o != "-b" || *a != "c" {
        tst.Error("failed to get_next4", o, a)
    }
}


func TestParseComplex(tst *testing.T) {
    t := func (tst *testing.T, x string) {
        c, e := parseComplex(x)
        s := fmt.Sprintf("%v", c)
        if s != x || e != "" {
            tst.Error("failed to TestParseComplex", c, e)
        }
    }
    t(tst, "(1.2+3.4i)")
    t(tst, "(1.2-3.4i)")
    t(tst, "(-1.2+3.4i)")
    t(tst, "(-1.2-3.4i)")
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


func TestParseArg1(tst *testing.T) {
    type A1 struct {
        B bool  `-b --bool |BOOL
                !test for bool1`
        V bool  `-v --verbose |VERBOSE
                test for bool2`
    }

    a1 := A1{}
    msg, args := ParseArg(&a1, []string{"-b"})
    if ! a1.B || a1.V || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    a2 := A1{}
    msg, args = ParseArg(&a2, []string{"-v"})
    if a2.B || ! a2.V || msg != "Missed option -b" {
        tst.Error("failed testParseArg 1", a2, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-bv"})
    if ! a1.B || ! a1.V || msg != "" {
        tst.Error("failed testParseArg 3", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 4", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-bv", "aa", "bb"})
    if ! a1.B || ! a1.V || msg != "" {
        tst.Error("failed testParseArg 3", a1, msg)
    }
    if len(args) != 2 || args[0] != "aa" || args[1] != "bb" {
        tst.Error("failed testParseArg 4", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-b", "-v"})
    if ! a1.B || ! a1.V || msg != "" {
        tst.Error("failed testParseArg 3", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 4", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-b", "-v", "11", "22"})
    if ! a1.B || ! a1.V || msg != "" {
        tst.Error("failed testParseArg 3", a1, msg)
    }
    if len(args) != 2 || args[0] != "11" || args[1] != "22" {
        tst.Error("failed testParseArg 4", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-bo", "-v", "11", "22"})
    if ! a1.B || a1.V || msg != "Don't know option: -o" {
        tst.Errorf("failed testParseArg 3 %v, [%s]", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 4", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"--bool", "-v", "11", "22"})
    if ! a1.B || ! a1.V || msg != "" {
        tst.Error("failed testParseArg 3", a1, msg)
    }
    if len(args) != 2 || args[0] != "11" || args[1] != "22" {
        tst.Error("failed testParseArg 4", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-b", "--verbose", "11", "22"})
    if ! a1.B || ! a1.V || msg != "" {
        tst.Error("failed testParseArg 3", a1, msg)
    }
    if len(args) != 2 || args[0] != "11" || args[1] != "22" {
        tst.Error("failed testParseArg 4", args)
    }
}


func TestParseArg2(tst *testing.T) {
    type A1 struct {
        I int  `-i --int 
                !test for int1`
        U uint `-u --uint 
                test for uint2`
    }

    a1 := A1{I:12}
    msg, args := ParseArg(&a1, []string{"-i"})
    if a1.I != 12 || a1.U != 0 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    a2 := A1{}
    msg, args = ParseArg(&a2, []string{"-i"})
    if a2.I != 0 || a2.U != 0 || msg != "option -i need parameter" {
        tst.Error("failed testParseArg 1", a2, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-u"})
    if a1.I != 0 || a1.U != 0 || msg != "option -u need parameter" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"--uint"})
    if a1.I != 0 || a1.U != 0 || msg != "option --uint need parameter" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    a4 := A1{}
    msg, args = ParseArg(&a4, []string{"-u", "32"})
    if a4.I != 0 || a4.U != 32 || msg != "Missed option -i" {
        tst.Error("failed testParseArg 1", a4, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{I: 12}
    msg, args = ParseArg(&a1, []string{"-iu", "32"})
    if a1.I != 12 || a1.U != 32 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{I: 12, U: 34}
    //msg, args = ParseArg(&a1, []string{"--int", "--uint"})
    msg, args = ParseArg(&a1, []string{"--int"})
    if a1.I != 12 || a1.U != 34 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-i", "-12"})
    if a1.I != -12 || a1.U != 0 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{I: -12}
    msg, args = ParseArg(&a1, []string{"-u", "-34"})
    if a1.I != -12 || a1.U != 0 || msg != "option -u need parameter" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-i", "-"})
    if a1.I != 0 || a1.U != 0 || msg != "option -i need parameter" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{I: -12}
    msg, args = ParseArg(&a1, []string{"-i", "-"})
    if a1.I != -12 || a1.U != 0 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 1 || args[0] != "-" {
        tst.Error("failed testParseArg 2", args)
    }
}


func TestParseArg3(tst *testing.T) {
    type A1 struct {
        Is []int    `-i --int 
                    !test for int1`
        Ss []string `-s --string
                    test for string2`
    }

    a1 := A1{}
    msg, args := ParseArg(&a1, []string{"-i", "12"})
    if len(a1.Is) != 1 || a1.Is[0] != 12 || len(a1.Ss) != 0 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{Is:[]int{12}}
    msg, args = ParseArg(&a1, []string{"-i", "-34", "-s", "-"})
    if a1.Is[0] != -34 || a1.Ss[0] != "-" || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }

    reset()
    a1 = A1{Is:[]int{12}}
    msg, args = ParseArg(&a1, []string{"-i", "-s", "-2"})
    if a1.Is[0] != 12 || len(a1.Ss) != 0 || msg != "option -s need parameter" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(a1.Is) != 1 {
        tst.Error("failed testParseArg 1", a1)
    }

    reset()
    a1 = A1{Is:[]int{12}}
    msg, args = ParseArg(&a1, []string{"-i", "-i", "-2"})
    if len(a1.Is) != 2 || a1.Is[0] != 12 || a1.Is[1] != -2 {
        tst.Error("failed testParseArg 1", a1, msg)
    }
}


func TestParseArg4(tst *testing.T) {
    type A1 struct {
        I int    `-i --int 
                    [1, 2, 3]
                    !test for int1`
        S string `-s --string
                    ["ab", "cd"]
                    test for string2`
    }

    a1 := A1{S: "ab"}
    msg, args := ParseArg(&a1, []string{"-i", "12"})
    if a1.I != 0 || a1.S != "ab" || msg != "option -i should be one of [1, 2, 3]" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = A1{}
    msg, args = ParseArg(&a1, []string{"-i", "2"})
    if a1.I != 2 || a1.S != "" || msg != "Missed option -s" {
        tst.Error("failed testParseArg 1", a1, msg)
    }

    reset()
    a1 = A1{S: "cd"}
    msg, args = ParseArg(&a1, []string{"-i", "2", "-s", "xy"})
    if a1.I != 2 || a1.S != "cd" || msg != `option -s should be one of ["ab", "cd"]` {
        tst.Error("failed testParseArg 1", a1, msg)
    }

    reset()
    a1 = A1{S: "cd"}
    msg, args = ParseArg(&a1, []string{"-i", "2", "-s", "ab"})
    if a1.I != 2 || a1.S != "ab" || msg != `` {
        tst.Error("failed testParseArg 1", a1, msg)
    }
}


func TestParseArg5(tst *testing.T) {
    type aa struct {
        I int    `-i --int 
                    [1, 2, 3]
                    test for int1`
    }

    a1 := aa{I: 2}
    msg, args := ParseArg(&a1, []string{})
    if a1.I != 2 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 2", args)
    }

    reset()
    a1 = aa{I: 1}
    msg, args = ParseArg(&a1, []string{})
    if a1.I != 1 || msg != "" {
        tst.Error("failed testParseArg 1", a1, msg)
    }

    reset()
    //a1 = aa{I: 2}
    a1 = aa{}
    msg, args = ParseArg(&a1, []string{"-i"})
    if msg != "option -i need parameter" {
        tst.Error("failed testParseArg 1", a1, msg)
    }
}


func TestParseGrp2(tst *testing.T) {
    type g2 struct {
        S string   `---FILENAME
                    !compress/decompress file
                    -c --compress
                    compress file
                    -x --decompress
                    decomrepss file`
    }

    g := g2{}
    msg, args := ParseArg(&g, []string{"abcd.txt"})
    if msg != "Missed option -c|-x" {
        tst.Error("failed testParseArg 1", g, msg)
    }
    if len(args) != 1 || args[0] != "abcd.txt" {
        tst.Error("failed testParseArg 1", args)
    }

    reset()
    g = g2{}
    msg, args = ParseArg(&g, []string{"-c", "abcd.txt"})
    if g.S != "c abcd.txt" || msg != "" {
        tst.Error("failed testParseArg 1", g, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 1", args)
    }

    reset()
    g = g2{}
    msg, args = ParseArg(&g, []string{"--decompress", "abcd.txt"})
    if g.S != "x abcd.txt" || msg != "" {
        tst.Error("failed testParseArg 1", g, msg)
    }
}


func TestParseGrp3(tst *testing.T) {
    type g3 struct {
        S string   `---FILENAME ""|
                    !compress/decompress file
                    -c --compress
                    compress file
                    -x --decompress
                    decomrepss file`
    }

    g := g3{}
    msg, args := ParseArg(&g, []string{"abcd.txt"})
    if msg != "Missed option -c|-x" {
        tst.Error("failed testParseArg 1", g, msg)
    }
    if len(args) != 1 || args[0] != "abcd.txt" {
        tst.Error("failed testParseArg 1", args)
    }


    reset()
    g = g3{}
    msg, args = ParseArg(&g, []string{"-c"})
    //tst.Log(msg, args)
    if msg != "" || len(args) != 0  {
        tst.Error("failed testParseGrp3", g, msg, args)
    }
}

/*
func TestParseGrp4(tst *testing.T) {
    type grp struct {
        Sel string
        Val string
    }
    type g3 struct {
        S *grp      `---FILENAME
                    compress/decompress file
                    -c --compress
                    compress file
                    -x --decompress
                    decomrepss file`
    }

    g := g3{S: &grp{Sel: "c", Val:"1234.txt"}}
    msg, args := ParseArg(&g, []string{"abcd.txt"})
    if g.S.Sel != "c" || g.S.Val != "1234.txt" || msg != "" {
        tst.Error("failed testParseArg 1", g, msg)
    }
    if len(args) != 1 || args[0] != "abcd.txt" {
        tst.Error("failed testParseArg 1", args)
    }

    reset()
    g = g3{S: &grp{Sel: "c", Val:"1234.txt"}}
    msg, args = ParseArg(&g, []string{"-x", "abcd.txt"})
    if g.S.Sel != "x" || g.S.Val != "abcd.txt" || msg != "" {
        tst.Error("failed testParseArg 1", g, msg)
    }
    if len(args) != 0 {
        tst.Error("failed testParseArg 1", args)
    }
}
*/

// go test -bench=. coap
func BenchmarkAdd(bm *testing.B) {
    src := []byte("52")
    var n byte
    for i := 0; i < bm.N; i++ {
        for _, b := range src {
            if b < '0' || b > '9' { break }
            n = n * 10 + b - '0'
        }
    }
}


func BenchmarkConv(bm *testing.B) {
    src := []byte("52")
    var n uint64
    for i := 0; i < bm.N; i++ {
        n, _ = strconv.ParseUint(string(src), 10, 32)
        _ = uint8(n)
    }
}


func BenchmarkN2S(bm *testing.B) {
    n := 1234
    for i := 0; i < bm.N; i++ {
        _ = fmt.Sprintf("%d", n)
    }
}


func BenchmarkS2N(bm *testing.B) {
    src := "1234"
    for i := 0; i < bm.N; i++ {
        _, _ = strconv.ParseInt(src, 10, 64)
    }
}


