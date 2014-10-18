package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "fmt"
)

func main () {
    // TODO remove
    var _ gocache.Cache
    var _ = fmt.Println
    //
    db := gocache.MakeCache()
    fmt.Println(db)
    db.Set("hi", "pop")
    fmt.Println(db)
    fmt.Println(db.Get("hi"))
    fmt.Println(db)
    fmt.Println(db.Delete("hi"))
    fmt.Println(db)
}
