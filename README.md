# coap
Command line Option and Argument Parser

```go

type Grp struct {
    sel string  // first must be string
    value int   // value is whatever can conver from []byte
}

type Arg struct {
    // Optional arguments

    // first line is short or long options,
    // second line if not start with {, will be help 
    // boolean option, only short
    B1 bool   `-v
               increase output verbosity`
    // boolean option, short and long, sep by space
    B2 bool   `-v --verbose
               increase output verbosity`

    // option, short and long, sep by space
    S1 string  `-n --name
                user name`
    // option, limit the selection
    S2 string  `-t --type
                {"admin", "worker"}
                user type`
    // option, you can put default, if no dfault,
    // then this arg must exists
    S3 string  `-n --name default
                user name`

    // group, can select one of them, but no more then one
    // if entry os string, result will combine with option and value
    // if entry is Grp, will fill first two fields

    // help will have multi line, follow same order
    C1 string  `-b|-e --begin|--end
                begin the service
                end the service`

    // group, can select one of them, but no more then one
    // default has two part, default option and value
    C2 string  `-c|-d --compress|--decompress '-d'|"filename"
                compress file 
                decompress file`

    // group, can select one of them, but no more then one
    // no default option, but has default value,
    // means must has one of these options, but followed value can omit
    C3 *Grp    `-c|-d --compress|--decompress ''|"filename"
                compress file 
                decompress file`
}
```
