package coap

import (
    "io"
    "os"
    "fmt"
    "bytes"
    "unicode"
    "strings"
    "reflect"
    "encoding/json"
)


type oaItem struct {    // option, argument item
    Short   string      // short name
    Vname   string      // name for value in help
    Long    string      // long name
    Must    bool        // must exists
    HasDft  bool        // has default
    MsgDft  string      // mseeage in help about default value
    StrDft  string      // str(default)
    Got     bool        // this item got from command line
    HelpLs  []string    // help message lines
    Cans    []string    // candidates
    Grp     []*oaItem   // this is a group item if not nil
    val     reflect.Value
    rsf     reflect.StructField
}


type oaInfo struct {
    oas []*oaItem
}


var infos map[uintptr]*oaInfo


func init() {
    infos = make(map[uintptr]*oaInfo)
}


func splitSpaceF(line string, doneF func ([]string) bool) (ret []string) {
    ret = make([]string, 0, 2)
    bg := 0
    for idx, r := range line {
        if ! unicode.IsSpace(r) { continue }
        if idx > bg {
            if len(ret) > 0 && doneF != nil && doneF(ret) {
                ret = append(ret, line[bg:])
                return
            }
            ret = append(ret, line[bg:idx])
        }
        bg = idx + 1
    }
    if bg < len(line) {
        ret = append(ret, line[bg:])
    }
    return
}


func (oa *oaItem)splitOpt(line string) (ret []string) {
    ret = splitSpaceF(line, func(r []string) bool {
                return strings.HasPrefix(r[len(r) - 1], "--")
           })
    return
}


func isZero(v reflect.Value) (b bool) {
    switch v.Kind() {
    // case reflect.Bool:
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
         return v.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
         return v.Uint() == 0
    case reflect.Float32, reflect.Float64:
         return v.Float() == 0.0
    //case reflect.Complex64, reflect.Complex128:
    //    return v.Complex() 
    case reflect.Slice, reflect.String:
        return v.Len() == 0
    default:
        panic("Not support type " + v.Kind().String())
    }
}


func (oa *oaItem)initDefault(val reflect.Value, dat string) {
    ds := strings.SplitN(dat, "|", 2)
    if len(ds) == 2 { oa.MsgDft = ds[1] }
    if ds[0] != "" {
        oa.HasDft = true
        oa.StrDft = ds[0]
    } else if ! isZero(val) {
        oa.HasDft = true
    }
}


func (oa *oaItem)initOpts(line string) string {
    // 1: short, 2: for long, 3: default
    opts := oa.splitOpt(line)
    for _, opt := range opts {
        fmt.Printf("[%s]\n", opt)
        switch {
        case opt[:2] == "--":
            oa.Long = opt[2:]
        case opt[:1] == "-":
            oa.Short = opt[1:2]
            if len(opt) > 2 {
                oa.Vname = opt[2:]
            }
        default:
            return opt
        }
    }
    return ""
}


func (oa *oaItem)initCans(line []byte) {
    vs := make([]interface{}, 0, 5)
    err := json.Unmarshal(line, &vs)
    if err != nil { panic("Failed to process candicates," + err.Error()) }
    oa.Cans = make([]string, 0, 5)
    for _, v := range vs {
        oa.Cans = append(oa.Cans, fmt.Sprint(v))
    }
}


func (oa *oaItem)initHelp(line string) {
    if oa.HelpLs == nil {
        //if line == "" { return }
        oa.HelpLs = make([]string, 0, 10)
        if strings.HasPrefix(line, "!") {
            oa.Must = ! strings.HasPrefix(line, "!!")
            oa.HelpLs = append(oa.HelpLs, strings.TrimSpace(line[1:]))
            return
        }
    }
    //oa.HelpLs = append(oa.HelpLs, strings.TrimSpace(line))
    oa.HelpLs = append(oa.HelpLs, line)
}


func (oa *oaItem)init(rsf reflect.StructField, val reflect.Value) {
    oa.rsf = rsf
    oa.val = val
    //tagLines := strings.Split(rsf.Tag, "\n")
    //fmt.Println("tags =", tagLines)
    isGrp := false
    soa := oa
    for _, l := range strings.Split(string(rsf.Tag), "\n") {
        line := strings.TrimSpace(l)
        switch {
        case strings.HasPrefix(line, "---"):    // group
            isGrp = true
            oa.Grp = make([]*oaItem, 0, 5)
            ret := splitSpaceF(line, func(r []string) bool { return len(r) > 0 })
            oa.Vname = ret[0][3:]
            if len(ret) > 1 { oa.initDefault(val, ret[1]) }
        case strings.HasPrefix(line, "-"):      // short or long
            if isGrp {
                soa = &oaItem{}
                oa.Grp = append(oa.Grp, soa)
            }
            dft := soa.initOpts(line)
            if ! isGrp {
                soa.initDefault(val, dft)
            }
        case strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]"):
            // candidates
            soa.initCans([]byte(line))
        default:    // help msg
            soa.initHelp(line)
        }
    }
}


func (oa *oaItem)parse(args []string) (cnt int) {
    // parse option/argument, return how many args used by this item
    return
}


func (oa *oaItem)helpShort(w io.Writer) {
    if oa.Short == "" && oa.Long == "" && len(oa.Grp) == 0 {
        return
    }
    if ! oa.Must { fmt.Fprintf(w, "[") }
    if len(oa.Grp) > 0 {    // group entry
        s := ""
        for _, g := range oa.Grp {
            fmt.Fprintf(w, s)
            g.helpShort(w)
            s = "|"
        }
    } else {                // regular entry
        if oa.Short != "" {
            fmt.Fprintf(w, "-%s", oa.Short)
        } else if oa.Long != "" {
            fmt.Fprintf(w, "--%s", oa.Long)
        }
    }
    if oa.Vname != "" {
        fmt.Fprintf(w, " ")
        if oa.HasDft { fmt.Fprintf(w, "[") }
        fmt.Fprintf(w, oa.Vname)
        if oa.HasDft { fmt.Fprintf(w, "]") }
    }
    if ! oa.Must { fmt.Fprintf(w, "]") }
}


func (oa *oaItem)helpLong(w io.Writer, align int) {
    s := bytes.Repeat([]byte(" "), align)
    b := bytes.Repeat([]byte(" "), align)
    copy(b[1:], "-" + oa.Short)
    copy(b[4:], "--" + oa.Long)
    w.Write(b)
    if len(oa.HelpLs) <= 0 { return }
    fmt.Fprintf(w, "%s\n", oa.HelpLs[0])
    for _, l := range oa.HelpLs {
        w.Write(s)
        fmt.Fprintf(w, "%s\n", l)
    }
}


type GRP struct {       // group
    sel     string
    val     string
}


type COAP struct {
    items   []oaItem
}


func verifySP(i interface{}) {  // Struct Pointer
    v := reflect.ValueOf(i)
    fmt.Println("v =", v)
    k := v.Kind()
    fmt.Println("k =", k)
    if k != reflect.Ptr {
        fmt.Fprintf(os.Stderr, "Need to be a ptr\n")
        os.Exit(1)
    }
    s := reflect.Indirect(v)
    fmt.Println("s =", s)
    k = s.Kind()
    fmt.Println("k =", k)
    if k != reflect.Struct {
        fmt.Fprintf(os.Stderr, "Need to be a struct\n")
        os.Exit(1)
    }
    a := s.Addr()
    fmt.Printf("a = %d\n", a)
}


func initial(i interface{}) *oaInfo {
    verifySP(i)
    v := reflect.ValueOf(i)
    info, ok := infos[v.Pointer()]
    if ok {     // already init, just return it
        return info
    }
    info = &oaInfo{oas: make([]*oaItem, 0, 5)}
    ii := reflect.Indirect(v)
    fmt.Println("ii =", ii)
    st := ii.Type()
    fmt.Println("st =", st)
    for idx := 0; idx < st.NumField(); idx++ {
        fs := st.Field(idx)
        fv := ii.Field(idx)
        fmt.Println("fv =", fv, reflect.TypeOf(fv))
        fv.SetString("MyName")
        it := &oaItem{}
        it.init(fs, fv)
        info.oas = append(info.oas, it)
    }
    return info
}


func Parse(arg interface{}) { ParseArg(arg, os.Args[1:]) }
func ParseArg(i interface{}, args []string) {
    initial(i)
    //fmt.Println("Parse(arg interface{})", arg)
    //st := reflect.TypeOf(arg)
    //fmt.Println(st)
    //fmt.Println(st.NumField())
}


func Help(arg interface{}) { HelpMsg(arg, "") }
func HelpMsg(i interface{}, msg string) {
    initial(i)
}
