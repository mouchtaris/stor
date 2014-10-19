package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/parser"
    "github.com/mouchtaris/topcoder_gocache/parser/action"
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
    "fmt"
    "strings"
    "os"
    "io"
)

func main () {
    r := strings.NewReader("set hello\r\nyou\r\nset hi\r\nme\r\nget hello hi\r\ndelete hi\r\n")
    var _ *lex.Lexer = lex.NewLexer(r)
    var _ gocache.Server
    var _ io.Reader = os.Stdin

    setConsumer := make(chan command.Set, 20)
    errors := make(chan error, 20)
    joiner := make(chan uint32, 20)
    go func () {
        defer func() { joiner <- 1 }()

        for setcomm := range setConsumer {
            fmt.Printf("set: %s\n", setcomm)
        }
        for err := range errors {
            fmt.Printf("error: %s\n", err)
        }
    }()

    lexer := lex.NewLexer(r)
    parser := parser.NewParser(lexer)
    parser.RegisterHandler(action.NewSet(setConsumer))
    err := parser.Parse()
    for ; err == nil; err = parser.Parse() {
    }

    errors <- err
    close(setConsumer)
    close(errors)
    for i := 0; i < 1; i++ {
        <-joiner
    }
}
