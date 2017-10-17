package lib

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "io"
    "strconv"
)

// Param : include the parameters.
// fmt.Println("Param: ", param)
type Param struct {
	h bool
    PageType bool

    PageLen int
    Start, End int

    Input io.ReadCloser
    Output io.WriteCloser
    
    startPage, endPage, printDest string
}

var param Param

func init() {
    flag.BoolVar(&param.h, "h", false, "get help")
    
    flag.IntVar(&param.PageLen, "l", 72, "specify `n` lines per page")

    flag.BoolVar(&param.PageType, "f", false, "specify '\\p' per page")

    flag.StringVar(&param.startPage, "s", "", "requires a `start page` number")
    flag.StringVar(&param.endPage, "e", "", "requires a `end page` number")

    flag.StringVar(&param.printDest, "d", "", "requires a `printer destination`")

    flag.Usage = usage
}

func usage() {
    fmt.Fprintf(os.Stderr, `selpg version 0.0.1
Usage: %v -s start_page -e end_page [ -f | -l lines_per_page ] [ -d dest ] [ in_filename ]

Options:
`, os.Args[0])
    flag.PrintDefaults()
}

func check()  {
    var err error;
    if param.Start, err = strconv.Atoi(param.startPage); err != nil {
        // or usage()
        // or possibly use `log.Fatalf` instead of:
        fmt.Fprintf(os.Stderr, "Missing required -%v argument/flag\n", "s")
        flag.Usage()
        // the same exit code flag.Parse uses
        os.Exit(2)
    }
    if param.End, err = strconv.Atoi(param.endPage); err != nil {
        fmt.Fprintf(os.Stderr, "Missing required -%v argument/flag\n", "e")
        flag.Usage()
        os.Exit(2)
    }

    if param.Start < 1 {
        fmt.Fprintf(os.Stderr, "%v: invalid start page Num %v\n", os.Args[0], param.Start)
        os.Exit(4)
    }
    if param.End < param.Start {
        fmt.Fprintf(os.Stderr, "%v: invalid end page Num %v\n", os.Args[0], param.End)
        os.Exit(4)
    }

    if !param.PageType && param.PageLen < 1 {
        fmt.Fprintf(os.Stderr, "%v: invalid page length %v\n", os.Args[0], param.PageLen)
    }

    if flag.NArg() > 0 {
        param.Input, err = os.Open(flag.Arg(0))
        if err != nil {
            fmt.Fprintf(os.Stderr, "Read file %q error\n", flag.Arg(0))
            panic(err)
        }
    } else {
        param.Input = os.Stdin
    }

    if param.printDest != "" {
        // cmd := exec.Command("lp -d" + param.printDest)
        cmd := exec.Command(param.printDest)
        param.Output, err = cmd.StdinPipe()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Write file %q error\n", param.printDest)
            panic(err)
        }
        cmd.Start()
    } else {
        param.Output = os.Stdout
    }
}

// Get : the Command Line parameters
// @Return : the Command Line parameters
func Get() *Param {
    flag.Parse()

    if param.h {
        flag.Usage()
        os.Exit(0)
    }

    check()

    return &param;
}

// Free : check the error
// and close io
func Free(err error) {
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
    }
    param.Input.Close()
    param.Output.Close()
    fmt.Fprintf(os.Stderr, "%v: done.\n", os.Args[0])
}
