package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "fmt"
    "strings"
)

func main () {
    r := strings.NewReader("set asok asda")
    parser := gocache.NewParser(r)
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
