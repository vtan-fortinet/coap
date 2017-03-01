# coap
Command line Option and Argument Parser

```go

type Arg struct {
    // Optional arguments

    // first line is short or long options,
    // second line if not start with [, will be help 
    // boolean option, only short
    B1 bool   `-v
               increase output verbosity`
    // boolean option, short and long, sep by space
    B2 bool   `-v --verbose
               increase output verbosity`

    // support -vvv
    A1 []bool `-v
               increase output verbosity`

    // option, short and long, sep by space
    S1 string  `-n --name
                user name`
    // option help name, default will use long name upper
    S2 string  `-nNAME --name
                user name`
    // option, limit the selection, is a json array of number or string
    S3 string  `-t --type
                ["admin", "worker"]
                user type`

    // Set default when instance the struct.
    // If the default is same as golang default, can attach on first line
    // if default not exists, means option and argument will stay together
    S4 string  `-p --password ""|default_desc_in_help
                user password`

    // Config must exists by leading help with !
    // So if not set default, then this means -s must exists, and must has arg
    // Start help with '!!' if help do need start with ! but not must exists
    // use '! !' if if help do need start with ! and must exists
    S5 string `-s --start
               !start command`

    // group, can select one of them, but no more then one
    // entry must be string, result will combine with option and value

    // --- start a group, default select will set in instance
    // default arg can set in instance or here
    // --- follow a name will ask argument
    G1 string  `---GRP
                Help for this group
                -u --upload
                help for upload, upload a file
                -d --download
                help for download, download file
    // you will get G1 as "u filename" or "d filename"
    // --- flollow nothing, will need no argument
    G2 string  `---
                Help for this group
                -b --beging
                help for begin, 
                -e --end
                help for end
    // you will get G2 "b" or "e"
}
```
