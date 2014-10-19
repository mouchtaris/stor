package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/parser"
    "github.com/mouchtaris/topcoder_gocache/parser/action"
    "github.com/mouchtaris/topcoder_gocache/parser/lex"
    "github.com/mouchtaris/topcoder_gocache/command"
    "github.com/mouchtaris/topcoder_gocache/cache"
    "fmt"
    "strings"
    "os"
    "io"
)

func open (fname string) io.Reader {
    r, err := os.Open(fname)
    if err != nil {
        return nil
    }
    return r
}

var Inputs = []io.Reader {
    open("inputs/00.txt"),
    open("inputs/01.txt"),
    open("inputs/02.txt"),
}

func main () {
    r := strings.NewReader("set hello\r\nyou\r\nset hi\r\nme\r\nget hello hi\r\ndelete hi\r\n" +
        "stats\r\nquit\r\n")
    var _ *lex.Lexer = lex.NewLexer(r)
    var _ gocache.Server
    var _ io.Reader = os.Stdin
    cache := cache.NewCache()

    comConsumer := make(chan command.Command, 20)
    errors := make(chan error, 20)
    joiner := make(chan uint32, 20)
    go func () {
        defer func() { joiner <- 1 }()
        for comm := range comConsumer {
            comm.PerformOn(cache)
        }
    }()
    go func () {
        defer func() { joiner <- 1 }()
        for err := range errors {
            fmt.Printf("error: %s\n", err)
        }
    }()

    lexer := lex.NewLexer(r)
    parser := parser.NewParser(lexer)
    parser.RegisterHandler(action.NewSet(comConsumer))
    parser.RegisterHandler(action.NewGet(comConsumer))
    parser.RegisterHandler(action.NewDelete(comConsumer))
    parser.RegisterHandler(action.NewStats(comConsumer))
    parser.RegisterHandler(action.NewQuit(comConsumer))
    err := parser.Parse()
    for ; err == nil; err = parser.Parse() {
    }

    errors <- err
    close(comConsumer)
    close(errors)
    for i := 0; i < 2; i++ {
        <-joiner
    }
}
