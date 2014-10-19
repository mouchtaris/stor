pkgdir = _build/topcoder.com/mouchtaris

${pkgdir}/%.o: src/topcoder.com/mouchtaris/%/*.go
	gccgo -I _build -g -O0 -c -o $@ -pedantic -Wall -Wextra $^

dirs = \
	${pkgdir}/scs/cache \
	${pkgdir}/scs/lex \
	${pkgdir}/scs/parser/action \
	${pkgdir}/scs/parser \

all: gcc go
gcc: ${dirs} scs
clean:
	rm -rvf _build pkg bin lol
${dirs}:
	mkdir -pv $@
scs: \
     ${pkgdir}/scs/util.o \
     ${pkgdir}/scs/cache.o \
     ${pkgdir}/scs/command.o \
     ${pkgdir}/scs/lex.o \
     ${pkgdir}/scs/parser/action.o \
     ${pkgdir}/scs/parser.o \
     ${pkgdir}/scs.o \
     ${pkgdir}/scs/main.o \

	gccgo -o $@ $^
go:
	go install topcoder.com/mouchtaris/scs/main

_build/topcoder_gocache_test.o: _build/topcoder_gocache.o



