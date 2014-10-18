package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/command"
    "fmt"
)

func main () {
    // TODO remove
    var _ gocache.Cache
    var _ = fmt.Println
    var _ command.Set
    //
    db := gocache.MakeCache()
    fmt.Println(db)
    db.Set("hi", "pop")
    fmt.Println(db)
    fmt.Println(db.Get("hi"))
    fmt.Println(db)
    fmt.Println(db.Delete("hi"))
    fmt.Println(db)
    fmt.Println(command.Set { })
}
