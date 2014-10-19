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
    r := strings.NewReader("set hello\r\nyou\r\nset hi\r\nme\r\nget hello hi\r\ndelete hi\r\n" +
        "stats\r\nquit\r\n")
    var _ *lex.Lexer = lex.NewLexer(r)
    var _ gocache.Server
    var _ io.Reader = os.Stdin

    setConsumer := make(chan command.Set, 20)
    getConsumer := make(chan command.Get, 20)
    delConsumer := make(chan command.Delete, 20)
    sttConsumer := make(chan command.Stats, 20)
    qitConsumer := make(chan command.Quit, 20)
    errors := make(chan error, 20)
    joiner := make(chan uint32, 20)
    go func () {
        defer func() { joiner <- 1 }()
        for setcomm := range setConsumer {
            fmt.Printf("set: %s\n", setcomm)
        }
    }()
    go func () {
        defer func() { joiner <- 1 }()
        for getcomm := range getConsumer {
            fmt.Printf("get: %s\n", getcomm)
        }
    }()
    go func () {
        defer func() { joiner <- 1 }()
        for comm := range delConsumer {
            fmt.Printf("del: %s\n", comm)
        }
    }()
    go func () {
        defer func() { joiner <- 1 }()
        for comm := range sttConsumer {
            fmt.Printf("stats: %s\n", comm)
        }
    }()
    go func () {
        defer func() { joiner <- 1 }()
        for comm := range qitConsumer {
            fmt.Printf("quit: %s\n", comm)
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
    parser.RegisterHandler(action.NewSet(setConsumer))
    parser.RegisterHandler(action.NewGet(getConsumer))
    parser.RegisterHandler(action.NewDelete(delConsumer))
    parser.RegisterHandler(action.NewStats(sttConsumer))
    parser.RegisterHandler(action.NewQuit(qitConsumer))
    err := parser.Parse()
    for ; err == nil; err = parser.Parse() {
    }

    errors <- err
    close(setConsumer)
    close(getConsumer)
    close(delConsumer)
    close(sttConsumer)
    close(qitConsumer)
    close(errors)
    for i := 0; i < 6; i++ {
        <-joiner
    }
}
