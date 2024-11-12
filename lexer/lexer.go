// sorry for not making the code really clean
package lexer

import (
	"log"
	"os"
	"teapot/report"
)

type TokenType int
const (
	// GO ENUMS SUCK
	// h
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

	// Arithmetic operators
	TT_PLUS
	TT_MINUS
	TT_STAR
	TT_SLASH
	TT_PERCENT

	// Bitwise operators
	TT_BAND
	TT_BOR
	TT_BNOT
	TT_BXOR

	// Comparison operators
	TT_EQUAL
	TT_NEQUAL
	TT_GREATER
	TT_LESS
	TT_GREATER_EQ
	TT_LESS_EQ

	// Logical operators
	TT_LAND
	TT_LOR
	TT_LNOT

	// Assignment operators
	TT_EQ
	TT_PLUS_EQ
	TT_MINUS_EQ
	TT_STAR_EQ
	TT_SLASH_EQ
	TT_PERCENT_EQ

	TT_BAND_EQ
	TT_BOR_EQ
	TT_BNOT_EQ
	TT_BXOR_EQ

	TT_INCREMENT
	TT_DECREMENT
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
		"TT_PLUS",
		"TT_MINUS",
		"TT_STAR",
		"TT_SLASH",
		"TT_PERCENT",
		"TT_BAND",
		"TT_BOR",
		"TT_BNOT",
		"TT_BXOR",
		"TT_EQUAL",
		"TT_NEQUAL",
		"TT_GREATER",
		"TT_LESS",
		"TT_GREATER_EQ",
		"TT_LESS_EQ",
		"TT_LAND",
		"TT_LOR",
		"TT_EQ",
		"TT_PLUS_EQ",
		"TT_MINUS_EQ",
		"TT_STAR_EQ",
		"TT_SLASH_EQ",
		"TT_PERCENT_EQ",
		"TT_BAND_EQ",
		"TT_BOR_EQ",
		"TT_BNOT_EQ",
		"TT_BXOR_EQ",
		"TT_INCREMENT",
		"TT_DECREMENT",
	}[tt]
}

type Token struct {
	Typ TokenType
	Lexeme string
	Line int
	Char int
}

func notoken() Token {
	return Token {Typ: TT_NONE, Lexeme: ""}
}

var line int = 1
var char int = 0

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
	// arithmetic operators
	case '+': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_PLUS, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '+': {
				*iter++
				return Token {Typ: TT_INCREMENT, Lexeme: "", Line: line, Char: char}
			}
			case '=': {
				*iter++
				return Token {Typ: TT_PLUS_EQ, Lexeme: "", Line: line, Char: char}
			}
		default:
			return Token {Typ: TT_PLUS, Lexeme: "", Line: line, Char: char}
		}
	}
	case '-': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_MINUS, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '-': {
				*iter++
				return Token {Typ: TT_DECREMENT, Lexeme: "", Line: line, Char: char}
			}
			case '=': {
				*iter++
				return Token {Typ: TT_MINUS_EQ, Lexeme: "", Line: line, Char: char}
			}
		default:
			return Token {Typ: TT_MINUS, Lexeme: "", Line: line, Char: char}
		}
	}
	case '*': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_STAR, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_STAR_EQ, Lexeme: "", Line: line, Char: char}
			}
		default:
			return Token {Typ: TT_STAR, Lexeme: "", Line: line, Char: char}
		}
	}
	case '/': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_SLASH, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '/': {
				// comment
				for !is_at_end(*iter, dat) && dat[*iter] != '\n' {
					*iter++
				}
				char = 0
				line++
				return Token {Typ: TT_NONE, Lexeme: "", Line: line, Char: char}
			}
			case '=': {
				*iter++
				return Token {Typ: TT_SLASH_EQ, Lexeme: "", Line: line, Char: char}
			}
		default:
			return Token {Typ: TT_SLASH, Lexeme: "", Line: line, Char: char}
		}
	}
	case '%': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_PERCENT, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_PERCENT_EQ, Lexeme: "", Line: line, Char: char}
			}
		default:
			return Token {Typ: TT_PERCENT, Lexeme: "", Line: line, Char: char}
		}
	}
	case '&': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_BAND, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_BAND_EQ, Lexeme: "", Line: line, Char: char - 1}
			}
			case '&': {
				*iter++
				return Token {Typ: TT_LAND, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_BAND, Lexeme: "", Line: line, Char: char}
		}
	}
	case '|': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_BOR, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_BOR_EQ, Lexeme: "", Line: line, Char: char - 1}
			}
			case '|': {
				*iter++
				return Token {Typ: TT_LOR, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_BOR, Lexeme: "", Line: line, Char: char}
		}
	}
	case '~': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_BNOT, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_BNOT_EQ, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_BNOT, Lexeme: "", Line: line, Char: char}
		}
	}
	case '^': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_BXOR, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_BXOR_EQ, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_BXOR, Lexeme: "", Line: line, Char: char}
		}
	}
	case '=': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_EQ, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_EQUAL, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_EQ, Lexeme: "", Line: line, Char: char}
		}
	}
	case '!': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_LNOT, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_NEQUAL, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_LNOT, Lexeme: "", Line: line, Char: char}
		}
	}
	case '>': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_GREATER, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_GREATER_EQ, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_GREATER, Lexeme: "", Line: line, Char: char}
		}
	}
	case '<': {
		if is_at_end(*iter + 1, dat) {
			return Token {Typ: TT_LESS, Lexeme: "", Line: line, Char: char}
		}

		switch dat[*iter + 1] {
			case '=': {
				*iter++
				return Token {Typ: TT_LESS_EQ, Lexeme: "", Line: line, Char: char - 1}
			}
		default:
			return Token {Typ: TT_LESS, Lexeme: "", Line: line, Char: char}
		}
	}
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
