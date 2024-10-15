package lexer

import (
	"fmt"
	"log"
	"os"
	"teapot/report"
)

type TokenType int
const (
	// GO ENUMS SUCK
	TT_NONE = iota

	TT_INT
	TT_STRING
	TT_ANY
	TT_ANYTYPE
	TT_BOOL
	TT_CHAR
	TT_CONST
	TT_FLOAT64
	TT_FLOAT32
	TT_VOID

	TT_IDENTIFIER

	TT_ORP
	TT_CRP

	TT_OSP
	TT_CSP

	TT_OBR
	TT_CBR

	TT_STR
	TT_INTEGER
	TT_HEXNUM
	TT_OCTNUM
	TT_FLOATNUM

	TT_SEMICOLON
)

func TT_ToStr(tt TokenType) string {
	return []string {
		"TT_NONE",
		"TT_INT",
		"TT_STRING",
		"TT_ANY",
		"TT_ANYTYPE",
		"TT_BOOL",
		"TT_CHAR",
		"TT_CONST",
		"TT_FLOAT64",
		"TT_FLOAT32",
		"TT_VOID",
		"TT_IDENTIFIER",
		"TT_ORP",
		"TT_CRP",
		"TT_OSP",
		"TT_CSP",
		"TT_OBR",
		"TT_CBR",
		"TT_STR",
		"TT_INTEGER",
		"TT_HEXNUM",
		"TT_OCTNUM",
		"TT_FLOATNUM",
		"TT_SEMICOLON",
	}[tt]
}

type Token struct {
	typ TokenType
	lexeme string
	line int
	char int
}

var line int = 1
var char int = 0

func is_numeric(r rune) bool {
	return (r >= '0' && r <= '9')
}

func is_alpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_'
}

func is_hex(r rune) bool {
	return is_numeric(r) || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f')
}

func is_at_end(iter int, dat []byte) bool {
	return iter == len(dat)
}

func lex_string(iter *int, dat []byte) Token {
	*iter++

	var lexeme string

	for dat[*iter] != '"' {
		if is_at_end(*iter + 1, dat) {
			report.Errorf("In Lexer :: Unterminated string at line %d", line)
			break
		}

		lexeme += string(dat[*iter])
		char++;
		if (dat[*iter] == '\n') {
			line++;
			char = 0;
		}
	}

	*iter++

	return Token {typ: TT_STR, lexeme: lexeme, line: line, char: char}
}

func lex_identifier(iter *int, dat []byte) Token {
	var token Token

	for is_alpha(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		token.lexeme += string(dat[*iter]);
		*iter++;
		char++;
	}

	token.typ = TT_IDENTIFIER;
	token.line = line;
	token.char = char;

	switch token.lexeme {
		case "int":
			token.typ = TT_INT
		case "string":
			token.typ = TT_STRING
		case "any":
			token.typ = TT_ANY
		case "anytype":
			token.typ = TT_ANYTYPE
		case "bool":
			token.typ = TT_BOOL
		case "char":
			token.typ = TT_CHAR
		case "const":
			token.typ = TT_CONST
		case "float64":
			token.typ = TT_FLOAT64
		case "float32":
			token.typ = TT_FLOAT32
		case "void":
			token.typ = TT_VOID
	}

	return token;
}

func lex_octal(iter *int, dat[]byte) Token {
	var token Token

	for is_numeric(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		if dat[*iter] >= '8' {
			report.Errorf("In Lexer :: Digits above 7 are not allowed in octal numbers.\n")
		}

		token.lexeme += string(dat[*iter])
	}

	token.typ = TT_OCTNUM
	token.line = line
	token.char = char

	return token
}

func lex_hex(iter *int, dat []byte) Token {
	var token Token

	token.lexeme += "0x"

	for is_hex(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		token.lexeme += string(dat[*iter])
	}

	token.typ = TT_HEXNUM
	token.line = line
	token.char = char

	return token
}

func lex_number(iter *int, dat []byte) Token {
	var token Token

	for is_numeric(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		token.lexeme += string(dat[*iter])
		*iter++
		char++
	}

	token.line = line;
	token.char = char;

	if dat[*iter] != '.' {
		token.typ = TT_INTEGER
		return token
	}

	token.lexeme += "."

	*iter++;

	for is_numeric(rune(dat[*iter])) {
		token.lexeme += string(dat[*iter])
		*iter++
		char++
	}

	token.typ = TT_FLOATNUM

	return token
}

func notoken() Token {
	return Token {typ: TT_NONE, lexeme: ""}
}

func lex_character(iter *int, dat []byte) Token {
	switch dat[*iter] {
	case '(':
		return Token {typ: TT_ORP, lexeme: "", line: line, char: char}
	case ')':
		return Token {typ: TT_CRP, lexeme: "", line: line, char: char}
	case '[':
		return Token {typ: TT_OSP, lexeme: "", line: line, char: char}
	case ']':
		return Token {typ: TT_CSP, lexeme: "", line: line, char: char}
	case '{':
		return Token {typ: TT_OBR, lexeme: "", line: line, char: char}
	case '}':
		return Token {typ: TT_CBR, lexeme: "", line: line, char: char}
	case ';':
		return Token {typ: TT_SEMICOLON, lexeme: "", line: line, char: char}
	case '"':
		return lex_string(iter, dat)
	case '0': {
		if is_at_end(*iter + 1, dat) {
			return Token {typ: TT_INTEGER, lexeme: "0", line: line, char: char}
		}

		if dat[*iter + 1] == 'x' {
			*iter += 2

			return lex_hex(iter, dat)
		}

		return lex_octal(iter, dat)
	}
	case ' ', '\t':
		return notoken()
	case '\n', '\r':
		line++
		char = 0
		return notoken()
	default: {
		if is_numeric(rune(dat[*iter])) {
			return lex_number(iter, dat)
		} else if is_alpha(rune(dat[*iter])) {
			return lex_identifier(iter, dat)
		} else {
			report.Errorf("Invalid character '%c'.\n", dat[*iter])
			return notoken()
		}
	}
	}
}

func LexFile(file string) []Token {
	dat, err := os.ReadFile(file)

	if err != nil {
		log.Fatalf("Couldn't read file: %s\n", dat)
		return nil
	}

	var i int = 0

	var tokens []Token

	for i < len(dat) {
		fmt.Printf("Lexing %c\n", dat[i])
		tokens = append(tokens, lex_character(&i, dat))
		i++
		char++
	}

	return tokens
}
