_build/github.com/mouchtaris/%.o: src/github.com/mouchtaris/%/*.go
	gccgo -I _build -c -o $@ -pedantic -Wall -Wextra $^

pkgdir = _build/github.com/mouchtaris/topcoder_gocache

all: gcc go
gcc: ${pkgdir} lol
clean:
	rm -rvf _build pkg bin
${pkgdir}:
	mkdir -pv $@
lol: \
     ${pkgdir}/command.o \
     ${pkgdir}.o \
     ${pkgdir}_test.o
	gccgo -o $@ $^
go:
	go install github.com/mouchtaris/topcoder_gocache_test

_build/topcoder_gocache_test.o: _build/topcoder_gocache.o



