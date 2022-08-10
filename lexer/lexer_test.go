package lexer

import "testing"
import "github.com/sashakoshka/arf/file"
import "github.com/sashakoshka/arf/types"

func TestTokenizeAll (test *testing.T) {
	file, err := file.Open("../tests/lexer/all")
	if err != nil {
		test.Log(err)
		test.Fail()
		return
	}
	
	tokens, err := Tokenize(file)
	if err == nil {
		test.Log("Tokenize() should have returned an error")
		test.Fail()
		return
	}

	correct := []Token {
		Token { kind: TokenKindSeparator },
		Token { kind: TokenKindPermission, value: types.Permission {
			Internal: types.ModeRead,
			External: types.ModeWrite,
		}},
		Token { kind: TokenKindReturnDirection },
		Token { kind: TokenKindInt, value: -349820394 },
		Token { kind: TokenKindUInt, value: 932748397 },
		Token { kind: TokenKindFloat, value: 239485.37520 },
		Token { kind: TokenKindString, value: "hello world\n" },
		Token { kind: TokenKindRune, value: 'E' },
		Token { kind: TokenKindName, value: "helloWorld" },
		Token { kind: TokenKindColon },
		Token { kind: TokenKindDot },
		Token { kind: TokenKindLBracket },
		Token { kind: TokenKindRBracket },
		Token { kind: TokenKindLBrace },
		Token { kind: TokenKindRBrace },
		Token { kind: TokenKindPlus },
		Token { kind: TokenKindMinus },
		Token { kind: TokenKindIncrement },
		Token { kind: TokenKindDecrement },
		Token { kind: TokenKindAsterisk },
		Token { kind: TokenKindSlash },
		Token { kind: TokenKindAt },
		Token { kind: TokenKindExclamation },
		Token { kind: TokenKindPercent },
		Token { kind: TokenKindTilde },
		Token { kind: TokenKindLessThan },
		Token { kind: TokenKindLShift },
		Token { kind: TokenKindGreaterThan },
		Token { kind: TokenKindRShift },
		Token { kind: TokenKindBinaryOr },
		Token { kind: TokenKindLogicalOr },
		Token { kind: TokenKindBinaryAnd },
		Token { kind: TokenKindLogicalAnd },
	}

	if len(tokens) != len(correct) {
		test.Log("lexed", len(tokens), "tokens, want", len(correct))
		test.Fail()
		return
	}
	test.Log("token slice length match", len(tokens), "=", len(correct))

	for index, token := range tokens {
		if !token.Equals(correct[index]) {
			test.Log("token", index, "not equal")
			test.Fail()
			return
		}
	}
	test.Log("token slice content match")
}
