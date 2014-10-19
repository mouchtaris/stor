_build/github.com/mouchtaris/%.o: src/github.com/mouchtaris/%/*.go
	gccgo -I _build -g -O0 -c -o $@ -pedantic -Wall -Wextra $^

pkgdir = _build/github.com/mouchtaris/topcoder_gocache
dirs = \
	${pkgdir}/cache \
	${pkgdir}/parser \
	${pkgdir}/parser/action \
	${pkgdir}/lex \

all: gcc go
gcc: ${dirs} lol
clean:
	rm -rvf _build pkg bin lol
${dirs}:
	mkdir -pv $@
lol: \
     ${pkgdir}/util.o \
     ${pkgdir}/cache.o \
     ${pkgdir}/command.o \
     ${pkgdir}/lex.o \
     ${pkgdir}/parser/action.o \
     ${pkgdir}/parser.o \
     ${pkgdir}.o \
     ${pkgdir}_test.o
	gccgo -o $@ $^
go:
	go install topcoder.com/mouchtaris/scs/main

_build/topcoder_gocache_test.o: _build/topcoder_gocache.o



