package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/command"
    "fmt"
    "strings"
)

func main () {
    // TODO remove
    var _ gocache.Cache
    var _ = fmt.Println
    var _ command.Set
    //
    db := gocache.NewCache()
    fmt.Println(db)
    db.Set("hi", "pop")
    fmt.Println(db)
    fmt.Println(db.Get("hi"))
    fmt.Println(db)
    fmt.Println(db.Delete("hi"))
    fmt.Println(db)
    fmt.Println(command.Set { })

    r := strings.NewReader("set asok asadsadpd")
    parser := gocache.NewParser(r)
    pn := func () {
        tok, err := parser.NextToken()
        if err == nil {
            fmt.Println(string(tok))
        } else {
            fmt.Println(err)
        }
    }
    pn()
    pn()
    pn()

    fmt.Println(byte("ad  d"[3]))
    fmt.Printf("%T\n", "sda")
    fmt.Printf("%T\n", "sda"[1])
}
