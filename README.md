# coap
Command line Option and Argument Parser

```go

type Grp struct {
    sel string  // first must be string
    val int     // value is whatever can conver from []byte
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
    // option help name, default will use long name upper
    S2 string  `-nNAME --name
                user name`
    // option, limit the selection
    S3 string  `-t --type
                {"admin", "worker"}
                user type`

    // Set default when instance the struct.
    // If the default is same as golang default, then can attach on first line
    // if not exists default, means option and argument will stay together
    S4 string  `-p --password ""|default_desc_in_help
                user password`

    // Config must exists by lead help with !
    // So if not set default, then this means -s must exists, and must has arg
    // Start help with '!!' if help do need start with ! but not must exists
    // use '! !' if if help do need start with ! and must exists
    S5 string `-s --start
               !start command`

    // group, can select one of them, but no more then one
    // if entry os string, result will combine with option and value
    // if entry is Grp, will fill first two fields

    // --- start a group, default select will set in instance
    // default arg can set in instance or here
    G1 string  `---GRP -b|"filename"
                Help for this group
                -b --begin
                help for begin: begin the service
                -e --end
                end the service`

    // use Grp struct, must has sel and val, sl must be string,
    // val can be anything convatable from string
    G2 *Grp    `---FILENAME
                !Compress/uncompress file
                -c --compress
                compress file 
                -d --decompress
                decompress file`
}
```
