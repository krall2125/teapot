package lexer

import (
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
	TT_RETURN

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
		"TT_RETURN",
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
	Typ TokenType
	Lexeme string
	Line int
	Char int
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

	var token Token

	token.Typ = TT_STR;
	token.Char = char;
	token.Line = line;

	for dat[*iter] != '"' {
		if is_at_end(*iter + 1, dat) {
			report.Errorf("In Lexer :: Unterminated string at line %d", line)
			break
		}

		token.Lexeme += string(dat[*iter])
		char++;
		if (dat[*iter] == '\n') {
			line++;
			char = 0;
		}

		*iter++;
	}

	return token;
}

func lex_identifier(iter *int, dat []byte) Token {
	var token Token

	token.Char = char;
	token.Line = line;

	for is_alpha(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		token.Lexeme += string(dat[*iter]);
		*iter++;
		char++;
	}

	token.Typ = TT_IDENTIFIER;

	*iter--
	char--

	switch token.Lexeme {
		case "int":
			token.Typ = TT_INT
		case "string":
			token.Typ = TT_STRING
		case "any":
			token.Typ = TT_ANY
		case "anyType":
			token.Typ = TT_ANYTYPE
		case "bool":
			token.Typ = TT_BOOL
		case "char":
			token.Typ = TT_CHAR
		case "const":
			token.Typ = TT_CONST
		case "float64":
			token.Typ = TT_FLOAT64
		case "float32":
			token.Typ = TT_FLOAT32
		case "void":
			token.Typ = TT_VOID
		case "return":
			token.Typ = TT_RETURN
	}

	return token;
}

func lex_octal(iter *int, dat[]byte) Token {
	var token Token

	token.Char = char
	token.Line = line

	for is_numeric(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		if dat[*iter] >= '8' {
			report.Errorf("In Lexer :: Digits above 7 are not allowed in octal numbers.\n")
		}

		token.Lexeme += string(dat[*iter])
		*iter++;
	}

	token.Typ = TT_OCTNUM

	*iter--
	char--

	return token
}

func lex_hex(iter *int, dat []byte) Token {
	var token Token

	token.Lexeme += "0x"

	token.Char = char
	token.Line = line

	for is_hex(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		token.Lexeme += string(dat[*iter])
		*iter++
	}

	token.Typ = TT_HEXNUM

	*iter--
	char--

	return token
}

func lex_number(iter *int, dat []byte) Token {
	var token Token

	token.Char = char;
	token.Line = line;

	for is_numeric(rune(dat[*iter])) && !is_at_end(*iter, dat) {
		token.Lexeme += string(dat[*iter])
		*iter++
		char++
	}

	if dat[*iter] != '.' {
		token.Typ = TT_INTEGER
		*iter--
		char--
		return token
	}

	token.Lexeme += "."

	*iter++;

	for is_numeric(rune(dat[*iter])) {
		token.Lexeme += string(dat[*iter])
		*iter++
		char++
	}

	*iter--
	char--

	token.Typ = TT_FLOATNUM

	return token
}

func notoken() Token {
	return Token {Typ: TT_NONE, Lexeme: ""}
}

func lex_character(iter *int, dat []byte) Token {
	switch dat[*iter] {
	case '(':
		return Token {Typ: TT_ORP, Lexeme: "", Line: line, Char: char}
	case ')':
		return Token {Typ: TT_CRP, Lexeme: "", Line: line, Char: char}
	case '[':
		return Token {Typ: TT_OSP, Lexeme: "", Line: line, Char: char}
	case ']':
		return Token {Typ: TT_CSP, Lexeme: "", Line: line, Char: char}
	case '{':
		return Token {Typ: TT_OBR, Lexeme: "", Line: line, Char: char}
	case '}':
		return Token {Typ: TT_CBR, Lexeme: "", Line: line, Char: char}
	case ';':
		return Token {Typ: TT_SEMICOLON, Lexeme: "", Line: line, Char: char}
	case '"':
		return lex_string(iter, dat)
	case '0': {
		if is_at_end(*iter + 1, dat) || !is_numeric(rune(dat[*iter + 1])) {
			return Token {Typ: TT_INTEGER, Lexeme: "0", Line: line, Char: char}
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

func LexStr(code []uint8) []Token {
	var i int = 0

	var tokens []Token

	for i < len(code) {
		tokens = append(tokens, lex_character(&i, code))
		i++
		char++
	}

	return Filter(tokens, func (a Token) bool {
		return a.Typ != TT_NONE
	})
}

func LexFile(file string) []Token {
	dat, err := os.ReadFile(file)

	if err != nil {
		log.Fatalf("Couldn't read file: %s\n", dat)
		return nil
	}

	return LexStr(dat);
}

func Filter[E any](s []E, f func(E) bool) []E {
	s2 := make([]E, 0, len(s))
	for _, e := range s {
		if f(e) {
			 s2 = append(s2, e)
		}
	}
	return s2
}
