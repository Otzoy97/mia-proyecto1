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
%token<TEXT> opType
%token<TEXT> delete_
%token<TEXT> opDelete_
%token<TEXT> opRep_
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
        if c == eof {
            return eof
        }
        switch state {
            case 0:
                if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
                    //Se concatena el caracter a la cadena
                    stringRec = string(c)
                    state = 1
                }
                else if c >= '0' && c <= '9' {
                    numberRec = int(c) - '0'
                    state = 3
                }
                else if c == '-' {
                    stringRec = string(c)
                    state = 4
                }
                else if c == '\\' {
                    state = 8
                }
                else if c == '#' {
                    state = 11
                }
                else if c == '"' {
                    //Indica el inicio de una cadena
                    stringRec = ""
                    state = 2
                }
                else if c == ' ', c == '\t', c == '\n', c == '\r'{
                    state = 0
                } else {
                    log.Printf("unrecognized character %q", c)
                }
            case 1:
                if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
                    stringRec += string(c)
                } else {
                    if c != eof {
                        x.peek = c 
                    }
                    yylval.TEXT = stringRec
                    return str_
                }
            case 2:
                if c == '"'{
                    if c != eof {
                        x.peek = c 
                    }
                    yylval.TEXT = stringRec
                    return str_
                } else if c != '\n' {
                    stringRec += string(c)
                } else {
                    log.Printf("unclosed string")
                }
            case 3:
                if c >= '0' && c <= '9' {
                    numberRec *= 10
                    numberRec += int(c) - '0'
                } else {
                    if c != eof {
                        x.peek = c
                    }
                    yylval.NUM = numberRec
                    return num_int
                }
            case 4:
                if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
                    stringRec += string(c)
                    state = 5
                }
                else if c >= '0' && c <= '9' {
                    numberRec = (int(c) - '0') * -1
                    state = 7
                }
                else if c == '>' {
                    yylval.TEXT = '->'
                    return assig
                } else {
                    log.Printf("unrecognized character %q", c)
                }
            case 5:
                if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
                    stringRec += string(c)
                } else {
                    if c != eof {
                        x.peek = c
                    }
                    return x.reservedWords(stringRec)
                }
            case 7:
                if c >= '0' && c <= '9' {
                    numberRec *= 10
                    numberRec = (int(c) - '0') * -1
                } else {
                    if c != eof {
                        x.peek = c
                    }
                    yylval.NUM = numberRec
                    return num_int   
                }
            case 8:
                if c == '*' {
                    state = 9
                } else {
                    if c != eof {
                        x.peek = c
                    }
                    state = 0
                }
            case 9, 10:
                if c == '\n'{
                    stringRec = ""
                    numberRec = 0
                    state = 0
                }
            default: 
                log.Printf("unrecognized character %q", c)
        }
	}
}

func (x *exprLex) reservedWords(s string) int {
    yylval.TEXT = stringRec
    switch s {
        case "mkdisk":
            return mkdisk
        case "rmdisk":
            return rmdisk
        case "fdisk":
            return fdisk
        case "mount": 
            return mount
        case "unmount":
            return unmount
        case "exec": 
            return exec
        case "rep": 
            return rep
        case "-path":
            return path
        case "-size":
            return size
        case "-fit":
            return fit
        case "-unit":
            return unit
        case "-type":
            return type
        case "-delete":
            return delete_
        case "-name":
            return name
        case "-add":
            return add
        case "-id":
            return id
        case "bf":
            return opFit
        case "ff":
            return opFit
        case "wf":
            return opFit
        case "m":
            return opUnit
        case "k":
            return opUnit
        case "b":
            return opUnit
        case "p":
            return opType
        case "e":
            return opType
        case "l":
            return opType
        case "fast":
            return opDelete_
        case "full":
            return opDelete_
        case "mbr":
            return opRep_
        case "disk":
            return opRep_
        default:
            return str_
    }

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
