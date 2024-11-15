package lexer_test

import (
	"testing"
	"teapot/lexer"
	// "github.com/stretchr/testify/assert"
)

func TestLexerBatch(t *testing.T) {
	t.Log("Starting lexer test 1 - batch.");

	code := `int string any anytype bool char const float64 float32 void return
		 asdfghjkl ()[]{}"maybe mayb e"123456 0x1a3b4f 07753 3.456;+-*/%
		 &|~^ ==!=>< >=<= &&||! = += -= *= /=%= &=|= ~=^= ++ --`;

	tokens := lexer.LexStr([]uint8(code))

	for i := range len(tokens) {
		t.Logf("%s '%s' %d %d\n",
			lexer.TT_ToStr(tokens[i].Typ),
			tokens[i].Lexeme,
			tokens[i].Line,
			tokens[i].Char)
	}
}
