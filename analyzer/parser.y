%{
    package cmd

    import (
        "bufio"
        "bytes"
        "fmt"
        "io"
        "log"
        "math/big"
        "os"
        "unicode/utf8"
        "strconv"
    )

    var pRaiz []int

    type option struct{
        name string
        value interface{}
    }
%}

%union{
    TEXT interface{}
    NUM int
    CMD *command
    OPTLST map[string]interface{}
}

%token<NUM> num_int
%token<TEXT> str_
%token<TEXT> mkdisk
%token<TEXT> rmdisk
%token<TEXT> fdisk
%token<TEXT> mount
%token<TEXT> unmount
%token<TEXT> exec

%token<TEXT> path
%token<TEXT> size
%token<TEXT> fit
%token<TEXT> opFit
%token<TEXT> unit
%token<TEXT> opUnit
%token<TEXT> type
%token<TEXT> delete_
%token<TEXT> name
%token<TEXT> add
%token<TEXT> id
%token<TEXT> assig

%type <OPTLST> OPT
%type <TEXT> MKOP
%type <TEXT> CMDLIST
%type <command> CMD

%start INIT

%%

INIT: CMDLIST | 

CMDLIST: 
    CMDLIST CMD 
    { 
        pRaiz = append(pRaiz, $2) 
    }
|   CMD 
    { 
        pRaiz = append(pRaiz, $1) 
    }

CMD: 
    mkdisk MKOP 
    {
        $$ = mkdisk{}
    }

MKOP: 
    MKOP OPT 
    {
        $1[$2.(option).name] = $2.(option).value
    }
|   OPT 
    {
        $$ = make(map[string]interface{})
        $$[$1.(option).name] = $1.(option).value
    }

OPT:
    size assig num_int
    {
        $$ = option{name:"size", value: $3}
    }
|   fit assig opFit
    {
        $$ = option{name:"fit", value: $3}
    }
|   unit assig opUnit
    {
        $$ = option{name:"unit", value: $3}
    }
|   path assig str_
    {
        $$ = option{name:"path", value: $3}
    }

%%

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (x *exprLex) Lex(yylval *exprSymType) int {
    state := 0
    stringRec := ""
    numberRec := 0
	for {
		c := x.next()
        switch state {
            case 0:
                if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
                    //Se concatena el caracter a la cadena
                    stringRec += string(c)
                    state = 1
                }
                else if c >= '0' && c <= '9' {
                    numberRec = int(c) - '0'
                    state = 3
                }
                else if c == '-' {
                    state = 4
                }
                else if c == '\\' {
                    state = 5
                }
                else if c == '#' {
                    state = 6
                }
                else if c == '"' {
                    //Indica el inicio de una cadena
                    state = 2
                }
                else if c == ' ', c == '\t', c == '\n', c == '\r'{
                    state = 0
                } else {
                    log.Printf("unrecognized character %q", c)
                }
            case 1:
                if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
                    state = 1
                } else {
                    return 
                }
            case 2:
                if c == '"'{
                    //Se presume que la cadena ya acabó
                    //RETURN CADENA
                } else if c == '\n' {
                    //Error
                }
            case 3:
                if c >= '0' && c <= '9' {

                }
        }
        
		switch c {
		case eof:
			return eof
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.num(c, yylval)
		case '-'
			return 

		// Recognize Unicode multiplication and division
		// symbols, returning what the parser expects.
		case '×':
			return '*'
		case '÷':
			return '/'

		case ' ', '\t', '\n', '\r':
		default:
			log.Printf("unrecognized character %q", c)
		}
	}
}

// Lex a number.
func (x *exprLex) num(c rune, yylval *exprSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
	L: for {
		c = x.next()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E':
			add(&b, c)
		default:
			break L
		}
	}
	if c != eof {
		x.peek = c
	}
	yylval.num = &big.Rat{}
	_, ok := yylval.num.SetString(b.String())
	if !ok {
		log.Printf("bad number %q", b.String())
		return eof
	}
	return NUM
}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *exprLex) Error(s string) {
	log.Printf("parse error: %s", s)
}
