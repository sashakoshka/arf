package lexer

import "testing"
import "github.com/sashakoshka/arf/file"

func TestTokenizeAll (test *testing.T) {
	file, err := file.Open("tests/parser/all")
	if err != nil {
		test.Log(err)
		test.Fail()
	}
	
	tokens, err := Tokenize(file)
	if err != nil {
		test.Log(err)
		test.Fail()
	}

	correct := []Token {
		Token { kind: TokenKindSeparator, },
	}

	if len(tokens) != len(correct) {
		test.Log("lexed", tokens, "tokens, want", correct)
	}

	for index, token := range tokens {
		if !token.Equals(correct[index]) {
			test.Log("token", index, "not equal")
			test.Fail()
		}
	}
}
