package coap

import (
    "io"
    "os"
    "fmt"
    "path"
    "bytes"
    "unicode"
    "strings"
    "strconv"
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
    oam map[string]*oaItem
    oas []*oaItem
    sp  int     // length of help leading space
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
    case reflect.Bool:
         //return v.Bool() == false
        return false
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
         return v.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
         return v.Uint() == 0
    case reflect.Float32, reflect.Float64:
         return v.Float() == 0.0
    //case reflect.Complex64, reflect.Complex128:
    //    return v.Complex() 
    case reflect.Slice, reflect.String:
        //fmt.Printf("v.Len() = %d, [%s]\n", v.Len(), v.String())
        return v.Len() == 0
    default:
        panic("Not support type " + v.Kind().String())
    }
}


func (oa *oaItem)initDefault(val reflect.Value, dat string) {
    ds := strings.SplitN(dat, "|", 2)
    if len(ds) == 2 { oa.MsgDft = ds[1] }
    //fmt.Printf("ds = %d, [%s], %v\n", len(ds), ds[0], oa.HasDft)
    if ds[0] != "" {
        oa.HasDft = true
        oa.StrDft = ds[0]
        //fmt.Printf("StrDft=[%s]\n", oa.StrDft)
    } else if ! isZero(val) {
        oa.HasDft = true
        oa.StrDft = fmt.Sprintf("%v", val.Interface())
    }
}


func (oa *oaItem)initOpts(line string) string {
    // 1: short, 2: for long, 3: default
    opts := oa.splitOpt(line)
    for _, opt := range opts {
        //fmt.Printf("[%s]\n", opt)
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
        if oa.Vname == "" && len(oa.Long) > 0 {
            fmt.Fprintf(w, strings.ToUpper(oa.Long))
        } else {
            fmt.Fprintf(w, oa.Vname)
        }
        if oa.HasDft { fmt.Fprintf(w, "]") }
    }
    if ! oa.Must { fmt.Fprintf(w, "]") }
}


func (oa *oaItem)helpLong(w io.Writer, align int) {
    s := bytes.Repeat([]byte(" "), align)
    b := bytes.Repeat([]byte(" "), align)
    copy(b[1:], "-" + oa.Short + ",")
    copy(b[5:], "--" + oa.Long)
    w.Write(b)
    if len(oa.HelpLs) <= 0 { return }
    fmt.Fprintf(w, "%s\n", oa.HelpLs[0])
    for _, l := range oa.HelpLs[1:] {
        w.Write(s)
        fmt.Fprintf(w, "%s\n", l)
    }
}


func canUse(val *reflect.Value, org *string) bool {
    if org == nil || *org == "--" || strings.HasPrefix(*org, "--") {
        return false
    }
    if ! strings.HasPrefix(*org, "-") || *org == "-" {
        return true
    }
    k := val.Kind()
    if k == reflect.Slice {
        k = val.Type().Elem().Kind()
    }

    switch k {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        _, e := strconv.ParseInt(*org, 10, 64)
        return e == nil
    case reflect.Float32, reflect.Float64:
        _, e := strconv.ParseFloat(*org, 64)
        return e == nil
    }
    return false
}


func setValue(val *reflect.Value, dat string) (got int, err string) {
    got = 2
    switch val.Kind() {
    case reflect.Bool:
        got = 1
        val.SetBool(true)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        i, e := strconv.ParseInt(dat, 10, 64)
        if e != nil { return 0, e.Error() }
        val.SetInt(i)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        u, e := strconv.ParseUint(dat, 10, 64)
        if e != nil { return 0, e.Error() }
        val.SetUint(u)
        //return v.Uint() == 0
    case reflect.Float32, reflect.Float64:
        f, e := strconv.ParseFloat(dat, 64)
        if e != nil { return 0, e.Error() }
        val.SetFloat(f)
        //return v.Float() == 0.0
    //case reflect.Complex64, reflect.Complex128:
    //    return v.Complex() 
    case reflect.String:
        val.SetString(dat)
    //case reflect.Slice:
        //fmt.Printf("v.Len() = %d, [%s]\n", v.Len(), v.String())
        //return v.Len() == 0
    default:
        return 0, "Not support type " + val.Kind().String()
    }
    return
}


func (oa *oaItem)parse(opt, arg *string) (got int, err string) {
    // parse option/argument
    // got = 0 : not match
    // got = 1 : took opt, but not arg
    // got = 2 : took opt and arg
    // err: error message
    eqp := strings.HasPrefix(*opt, "--" + oa.Long + "=")
    if ("-" + oa.Short) != *opt && ("--" + oa.Long) != *opt  && ! eqp {
        //fmt.Printf("s=[%s], l=[%s], o=[%s]\n",
        //           ("-" + oa.Short), ("--" + oa.Long), *opt)
        return
    }
    //var e string
    //switch oa.val.Kind() {
    //case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    //    _, e = strconv.ParseInt(dat, 10, 64)
    //case reflect.Float32, reflect.Float64:
    //    f, e := strconv.ParseFloat(dat, 64)
    //}
    //if oa.val.Kind() == reflect.Bool {
    //    oa.val.SetBool(true)
    //    //fv.SetString("MyName")
    //    return
    //}
    var op, pa string
    var cu bool
    op = *opt
    if eqp {
        op = (*opt)[:len(oa.Long) + 2]
        pa = (*opt)[len(oa.Long) + 3:]
    //} else if arg != nil {
    } else if cu = canUse(&oa.val, arg); cu {
        pa = *arg
    } else if oa.HasDft {
        pa = oa.StrDft
    } else {
        //if oa.Short != "" {
        //    err = "option -" + oa.Short + " need parameter"
        //} else {
        //    err = "option --" + oa.Long + " need parameter"
        //}
        err = "option " + op + " need parameter"
        return
    }
    if oa.val.Kind() == reflect.Slice {
        /*
        a := make([]string, 0, 2)
        v := reflect.ValueOf(&a).Elem()
        t := reflect.TypeOf(a)
        e := t.Elem()
        n := reflect.New(e).Elem()
        n.SetString("abc")
        v.Set(reflect.Append(v, n))
        fmt.Println(a, n)
        */
        v := reflect.New(oa.val.Type().Elem()).Elem()
        got, err = setValue(&v, pa)
        if err != "" {
            //fmt.Printf("setValue return err=[%s]", err)
            return
        }
        oa.val.Set(reflect.Append(oa.val, v))
    } else {
        got, err = setValue(&oa.val, pa)
    }
    if err == "" && got > 0 {
        oa.Got = true
        if eqp || ! cu { got = 1 }
    }
    //fmt.Printf("got=%d, err=[%s]", got, err)
    return
}


func verifySP(i interface{}) {  // Struct Pointer
    v := reflect.ValueOf(i)
    k := v.Kind()
    if k != reflect.Ptr {
        fmt.Fprintf(os.Stderr, "Need to be a ptr\n")
        os.Exit(1)
    }
    s := reflect.Indirect(v)
    k = s.Kind()
    if k != reflect.Struct {
        fmt.Fprintf(os.Stderr, "Need to be a struct\n")
        os.Exit(1)
    }
}


func initial(i interface{}) *oaInfo {
    verifySP(i)
    v := reflect.ValueOf(i)
    info, ok := infos[v.Pointer()]
    if ok {     // already init, just return it
        return info
    }
    info = &oaInfo{oas: make([]*oaItem, 0, 5), oam: make(map[string]*oaItem, 5)}
    ii := reflect.Indirect(v)
    st := ii.Type()
    for idx := 0; idx < st.NumField(); idx++ {
        fs := st.Field(idx)
        fv := ii.Field(idx)
        it := &oaItem{}
        it.init(fs, fv)
        info.oas = append(info.oas, it)
        c := 0
        if len(it.Short) > 0 {
            info.oam["-" + it.Short] = it
            c = c + 2                   // for len("-S")
        }
        if len(it.Long) > 0 {
            info.oam["--" + it.Long] = it
            c = c + len(it.Long) + 2
            if it.Short != "" { c = c + 2 }   // for len(", ")
        }
        if c > info.sp { info.sp = c }
    }
    return info
}


func oasParse(oas []*oaItem, opt, arg *string) (got int, err string) {
    //fmt.Printf("opt=[%s], arg=[%s]\n", *opt, *arg)
    for _, oa := range oas {
        got, err = oa.parse(opt, arg)
        if got != 0 || err != "" {
            return
        }
    }
    if got == 0 {
        err = "Don't know option: " + *opt + err
    }
    return
}


func Parse(i interface{}) []string {
    msg, ps := ParseArg(i, os.Args[1:])
    if msg != "" {
        if len(os.Args) <= 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
            msg = ""
        }
        HelpMsg(i, msg)
        os.Exit(1)
    }
    return ps
}

func get_next(idx int, args []string) (o, a *string) {
    o = &(args[idx])
    //if idx < (len(args) - 1) && (! strings.HasPrefix(args[idx + 1], "-") || args[idx + 1] == "-") {
    if idx < (len(args) - 1) {
         a = &(args[idx + 1])
    }
    return
}
func ParseArg(i interface{}, args []string) (msg string, ps []string) {
    oi := initial(i)
    got := 0
    ps = make([]string, 0, 5)
    for idx := 0; idx < len(args); idx++ {
        switch {
        case args[idx] == "--":
            ps = append(ps, args[idx:]...)
            break
        case strings.HasPrefix(args[idx], "--"):
            o, a := get_next(idx, args)
            got, msg = oasParse(oi.oas, o, a)
            if msg != "" { return }
            if got > 1 { idx = idx + 1 }
        case strings.HasPrefix(args[idx], "-"):
            o, a := get_next(idx, args)
            opts := *o
            //println("opts1=", opts)
            for i := 1; i < len(opts) - 1; i++ {
                x := "-" + opts[i:i+1]
                //println("x=", x, "i=", i)
                got, msg = oasParse(oi.oas, &x, nil)
                if msg != "" { return }
            }
            opts = opts[:1] + opts[len(opts) - 1:]
            //println("opts2=", opts)
            got, msg = oasParse(oi.oas, &opts, a)
            if msg != "" { return }
            if got > 1 { idx = idx + 1 }
        default:
            ps = append(ps, args[idx])
        }
    }
    for _, oa := range oi.oas {
        if oa.Must && ! oa.Got {
            if oa.Short != "" {
                msg = "Missed option -" + oa.Short
            } else {
                msg = "Missed option --" + oa.Long
            }
        }
    }
    return
}


func Help(arg interface{}) { HelpMsg(arg, "") }
func HelpMsg(i interface{}, msg string) {
    oi := initial(i)
    a := 0
    if msg != "" {
        fmt.Fprintf(os.Stdout, "%s\n", msg)
    }
    fmt.Fprintf(os.Stdout, "Usage: %s ", path.Base(os.Args[0]))
    for _, oa := range oi.oas {
        oa.helpShort(os.Stdout)
        if i := len(oa.Short) + len(oa.Long) + 8; i > a {
            a = i
        }
    }
    fmt.Fprintf(os.Stdout, "\n")
    for _, oa := range oi.oas {
        oa.helpLong(os.Stdout, a)
    }
    fmt.Fprintf(os.Stdout, "\n")
}
