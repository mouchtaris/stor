package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/command"
    "github.com/mouchtaris/topcoder_gocache/parser"
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
    fmt.Println(parser.NextToken(r))
}
