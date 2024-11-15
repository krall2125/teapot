// sorry for not making the code really clean
package lexer

import (
	"log"
	"os"
	"teapot/report"
	"teapot/utils"
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

func NewToken(typ TokenType, lexeme string, line int, char int) Token {
	return Token {
		Typ: typ,
		Lexeme: lexeme,
		Line: line,
		Char: char,
	}
}

type Lexer struct {
	line int
	char int
	iter int
	dat []byte
}

func NewLexer(dat []byte) Lexer {
	return Lexer {
		line: 1,
		char: 0,
		iter: 0,
		dat: dat,
	}
}

func (lexer *Lexer) advance() {
	lexer.iter++
	lexer.char++

	if lexer.is_at_end() {
		return
	}

	if lexer.dat[lexer.iter] == '\n' {
		lexer.line++
		lexer.char = 0
		lexer.iter++
	}
}

func (lexer *Lexer) is_at_end() bool {
	return lexer.iter == len(lexer.dat)
}

func (lexer *Lexer) is_at_end_next() bool {
	return lexer.iter + 1 == len(lexer.dat)
}

func (lexer *Lexer) lex_string() Token {
	lexer.advance()

	var token Token = NewToken(TT_STR, "", lexer.line, lexer.char)

	for lexer.dat[lexer.iter] != '"' {
		if lexer.is_at_end() {
			report.Errorf("In Lexer :: Unterminated string at line %d", lexer.line)
			break
		}

		token.Lexeme += string(lexer.dat[lexer.iter])

		lexer.advance()
	}

	return token;
}

func (lexer *Lexer) lex_identifier() Token {
	var token Token = NewToken(TT_IDENTIFIER, "", lexer.line, lexer.char)

	for is_alpha(lexer.dat[lexer.iter]) && !lexer.is_at_end() {
		token.Lexeme += string(lexer.dat[lexer.iter])
		lexer.advance()
	}

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

func (lexer *Lexer) lex_octal() Token {
	var token Token = NewToken(TT_OCTNUM, "", lexer.line, lexer.char)

	for is_numeric(lexer.dat[lexer.iter]) && !lexer.is_at_end() {
		if lexer.dat[lexer.iter] >= '8' {
			report.Errorf("In Lexer :: Digits above 7 are not allowed in octal numbers.\n")
		}

		token.Lexeme += string(lexer.dat[lexer.iter])
		lexer.advance()
	}

	token.Typ = TT_OCTNUM

	lexer.iter--
	lexer.char--

	return token
}

func (lexer *Lexer) lex_hex() Token {
	var token Token = NewToken(TT_HEXNUM, "0x", lexer.line, lexer.char)

	for is_hex(lexer.dat[lexer.iter]) && !lexer.is_at_end() {
		token.Lexeme += string(lexer.dat[lexer.iter])
		lexer.advance()
	}

	lexer.iter--
	lexer.char--

	return token
}

func (lexer *Lexer) lex_number() Token {
	var token Token = NewToken(TT_INTEGER, "", lexer.line, lexer.char)

	for is_numeric(lexer.dat[lexer.iter]) && !lexer.is_at_end() {
		token.Lexeme += string(lexer.dat[lexer.iter])
		lexer.advance()
	}

	if lexer.dat[lexer.iter] != '.' {
		lexer.iter--
		lexer.char--
		return token
	}

	token.Lexeme += "."

	lexer.advance()

	for is_numeric(lexer.dat[lexer.iter]) {
		token.Lexeme += string(lexer.dat[lexer.iter])
		lexer.advance()
	}

	lexer.iter--
	lexer.char--

	token.Typ = TT_FLOATNUM

	return token
}

func (lexer *Lexer) lex_character() Token {
	switch lexer.dat[lexer.iter] {
	case '(':
		return NewToken(TT_ORP, "", lexer.line, lexer.char)
	case ')':
		return NewToken(TT_CRP, "", lexer.line, lexer.char)
	case '[':
		return NewToken(TT_OSP, "", lexer.line, lexer.char)
	case ']':
		return NewToken(TT_CSP, "", lexer.line, lexer.char)
	case '{':
		return NewToken(TT_OBR, "", lexer.line, lexer.char)
	case '}':
		return NewToken(TT_CBR, "", lexer.line, lexer.char)
	case ';':
		return NewToken(TT_SEMICOLON, "", lexer.line, lexer.char)
	// arithmetic operators
	case '+': {
		if lexer.is_at_end_next() {
			return NewToken(TT_PLUS, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '+': {
			lexer.advance()
			return NewToken(TT_INCREMENT, "", lexer.line, lexer.char)
		}
		case '=': {
			lexer.advance()
			return NewToken(TT_PLUS_EQ, "", lexer.line, lexer.char)
		}
		default:
			return NewToken(TT_PLUS, "", lexer.line, lexer.char)
		}
	}
	case '-': {
		if lexer.is_at_end_next() {
			return NewToken(TT_MINUS, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '-': {
			lexer.advance()
			return NewToken(TT_DECREMENT, "", lexer.line, lexer.char)
		}
		case '=': {
			lexer.advance()
			return NewToken(TT_MINUS_EQ, "", lexer.line, lexer.char)
		}
		default:
			return NewToken(TT_MINUS, "", lexer.line, lexer.char)
		}
	}
	case '*': {
		if lexer.is_at_end_next() {
			return NewToken(TT_STAR, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.iter++
			return NewToken(TT_STAR_EQ, "", lexer.line, lexer.char)
		}
		default:
			return NewToken(TT_STAR, "", lexer.line, lexer.char)
		}
	}
	case '/': {
		if lexer.is_at_end_next() {
			return NewToken(TT_SLASH, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '/': {
			// comment
			for !lexer.is_at_end() && lexer.dat[lexer.iter] != '\n' {
				lexer.advance()
			}
			return NewToken(TT_NONE, "", lexer.line, lexer.char)
		}
		case '=': {
			lexer.iter++
			return NewToken(TT_SLASH_EQ, "", lexer.line, lexer.char)
		}
		default:
			return NewToken(TT_SLASH, "", lexer.line, lexer.char)
		}
	}
	case '%': {
		if lexer.is_at_end_next() {
			return NewToken(TT_PERCENT, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_PERCENT_EQ, "", lexer.line, lexer.char)
		}
		default:
			return NewToken(TT_PERCENT, "", lexer.line, lexer.char)
		}
	}
	case '&': {
		if lexer.is_at_end_next() {
			return NewToken(TT_BAND, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_BAND_EQ, "", lexer.line, lexer.char - 1)
		}
		case '&': {
			lexer.advance()
			return NewToken(TT_LAND, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_BAND, "", lexer.line, lexer.char)
		}
	}
	case '|': {
		if lexer.is_at_end_next() {
			return NewToken(TT_BOR, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_BOR_EQ, "", lexer.line, lexer.char - 1)
		}
		case '|': {
			lexer.advance()
			return NewToken(TT_LOR, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_BOR, "", lexer.line, lexer.char)
		}
	}
	case '~': {
		if lexer.is_at_end_next() {
			return NewToken(TT_BNOT, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_BNOT_EQ, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_BNOT, "", lexer.line, lexer.char)
		}
	}
	case '^': {
		if lexer.is_at_end_next() {
			return NewToken(TT_BXOR, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_BXOR_EQ, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_BXOR, "", lexer.line, lexer.char)
		}
	}
	case '=': {
		if lexer.is_at_end_next() {
			return NewToken(TT_EQ, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_EQUAL, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_EQ, "", lexer.line, lexer.char)
		}
	}
	case '!': {
		if lexer.is_at_end_next() {
			return NewToken(TT_LNOT, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_NEQUAL, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_LNOT, "", lexer.line, lexer.char)
		}
	}
	case '>': {
		if lexer.is_at_end_next() {
			return NewToken(TT_GREATER, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_GREATER_EQ, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_GREATER, "", lexer.line, lexer.char)
		}
	}
	case '<': {
		if lexer.is_at_end_next() {
			return NewToken(TT_LESS, "", lexer.line, lexer.char)
		}

		switch lexer.dat[lexer.iter + 1] {
		case '=': {
			lexer.advance()
			return NewToken(TT_LESS_EQ, "", lexer.line, lexer.char - 1)
		}
		default:
			return NewToken(TT_LESS, "", lexer.line, lexer.char)
		}
	}
	case '"':
		return lexer.lex_string()
	case '0': {
		if lexer.is_at_end_next() || !is_numeric(lexer.dat[lexer.iter + 1]) {
			return NewToken(TT_INTEGER, "0", lexer.line, lexer.char)
		}

		if lexer.dat[lexer.iter + 1] == 'x' {
			lexer.advance()
			lexer.advance()

			return lexer.lex_hex()
		}

		return lexer.lex_octal()
	}
	case ' ', '\t':
		return NewToken(TT_NONE, "", lexer.line, lexer.char)
	default: {
		if is_numeric(lexer.dat[lexer.iter]) {
			return lexer.lex_number()
		} else if is_alpha(lexer.dat[lexer.iter]) {
			return lexer.lex_identifier()
		} else {
			report.Errorf("Invalid character '%c'.\n", lexer.dat[lexer.iter])
			return NewToken(TT_NONE, "", lexer.line, lexer.char)
		}
	}
	}
}

func is_numeric(r uint8) bool {
	return (r >= '0' && r <= '9')
}

func is_alpha(r uint8) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_'
}

func is_hex(r uint8) bool {
	return is_numeric(r) || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f')
}

func LexStr(code []uint8) []Token {
	var tokens []Token

	var lexer Lexer = NewLexer(code)

	for !lexer.is_at_end() {
		tokens = append(tokens, lexer.lex_character())
		lexer.advance()
	}

	return utils.Filter(tokens, func (a Token) bool {
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
