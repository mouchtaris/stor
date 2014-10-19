package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/parser"
    "fmt"
    "strings"
    "os"
)

func main () {
    r := strings.NewReader("set asok asda")
    var _ *parser.Lexer = parser.NewLexer(r)
    var _ gocache.Server

    parser := parser.NewLexer(os.Stdin)
    pn := func () error {
        err := parser.ReadCommand()
        if err == nil {
            fmt.Println("toen: ", string(parser.Token()))
        } else {
            fmt.Println("error: ", err)
        }
        return err
    }
    for pn() == nil {}
}
