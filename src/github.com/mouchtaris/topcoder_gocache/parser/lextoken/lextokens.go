package lextoken

const (
    SET = iota
    GET
    DELETE
    STATS
    QUIT
    KEY
    VALUE
)

type Token uint8
