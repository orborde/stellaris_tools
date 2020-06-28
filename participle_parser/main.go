package main

import (
	"io"
	"log"
	"os"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"

	"github.com/alecthomas/participle"
)

type ConfigFile struct {
	Assignments []Assignment `@@`
}

type Assignment struct {
	Key   Expression `@@ "="`
	Value Expression `@@`
}

type Expression struct {
	String       *string        `@String`
	Atom         *string        `| @Ident`
	KeyValueList []KeyValueList `"{" @@ "}"`
}

type KeyValueList struct {
	Key   Expression  `@@ "="`
	Value *Expression `@@`
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	f := func() io.Reader {
		if len(os.Args) == 2 {
			f, err := os.Open(os.Args[1])
			must(err)
			return f
		} else {
			return os.Stdin
		}
	}()

	lex := lexer.Must(ebnf.New(`
    Comment = ("#" | "//") { "\u0000"…"\uffff"-"\n" } .
	Ident = (alpha | "_") { "_" | alpha | digit | ":"} .
	Number = ("-" | "." | digit) {"-" | "." | digit} .
	Whitespace = " " | "\t" | "\n" | "\r" .
	Punct = "=" | "{" | "}" .
	String = "\"" { ( "\u0000"…"\uffff"-"\""-"\\" | "\\" "\u0000"…"\uffff") } "\"" .

    alpha = "a"…"z" | "A"…"Z" .
    digit = "0"…"9" .
`))
	parser := participle.MustBuild(&ConfigFile{},
		participle.Lexer(lex))
	var configFile ConfigFile
	must(parser.Parse(f, &configFile))
}