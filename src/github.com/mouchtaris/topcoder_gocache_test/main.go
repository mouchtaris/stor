package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "fmt"
    "strings"
    "os"
)

func main () {
    r := strings.NewReader("set asok asda")
    var _ *gocache.Parser = gocache.NewParser(r)

    parser := gocache.NewParser(os.Stdin)
    pn := func () error {
        err := parser.ReadCommand()
        if err == nil {
            fmt.Println("toen: ", string(parser.Token()))
        } else {
            fmt.Println("error: ", err)
        }
        return err
    }
    pn()
    pn()
    pn()
    pn()
    pn()
    pn()
    pn()
    pn()
}
