package main

import "github.com/LinJinghua/selpg/lib"

func main() {
    param := lib.Get()
    
    var err error
    if param.PageType {
        err = lib.CopyPages(param.Input, param.Output, param.Start - 1, param.End)
    } else {
        err = lib.CopyLines(param.Input, param.Output, (param.Start - 1) * param.PageLen, param.End * param.PageLen)
    }
    lib.Free(err)
}
