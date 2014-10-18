_build/github.com/mouchtaris/%.o: src/github.com/mouchtaris/%/*.go
	gccgo -I _build -c -o $@ -pedantic -Wall -Wextra $^

all: gcc go
gcc: _build/github.com/mouchtaris/ lol
clean:
	rm -rvf _build pkg bin
_build/github.com/mouchtaris/:
	mkdir -pv $@
lol: _build/github.com/mouchtaris/topcoder_gocache.o \
     _build/github.com/mouchtaris/topcoder_gocache_test.o
	gccgo -o $@ $^
go:
	go install github.com/mouchtaris/topcoder_gocache_test

_build/topcoder_gocache_test.o: _build/topcoder_gocache.o



