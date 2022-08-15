package lexer

import "testing"
import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/types"

func checkTokenSlice (filePath string, correct []Token, test *testing.T) {
	file, err := file.Open(filePath)
	if err != nil {
		test.Log(err)
		test.Fail()
		return
	}
	
	tokens, err := Tokenize(file)

	// print all tokens
	for index, token := range tokens {
		test.Log(index, "\tgot token:", token.Describe())
	}
	
	if err != nil {
		test.Log("returned error:")
		test.Log(err.Error())
		test.Fail()
		return
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
			test.Log (
				"have", token.Describe(),
				"want", correct[index].Describe())
			test.Fail()
			return
		}
	}
	test.Log("token slice content match")
}

func TestTokenizeAll (test *testing.T) {
	checkTokenSlice("../tests/lexer/all.arf", []Token {
		Token { kind: TokenKindSeparator },
		Token { kind: TokenKindPermission, value: types.Permission {
			Internal: types.ModeRead,
			External: types.ModeWrite,
		}},
		Token { kind: TokenKindReturnDirection },
		Token { kind: TokenKindInt, value: int64(-349820394) },
		Token { kind: TokenKindUInt, value: uint64(932748397) },
		Token { kind: TokenKindFloat, value: 239485.37520 },
		Token { kind: TokenKindString, value: "hello world!\n" },
		Token { kind: TokenKindRune, value: 'E' },
		Token { kind: TokenKindName, value: "helloWorld" },
		Token { kind: TokenKindColon },
		Token { kind: TokenKindDot },
		Token { kind: TokenKindComma },
		Token { kind: TokenKindLBracket },
		Token { kind: TokenKindRBracket },
		Token { kind: TokenKindLBrace },
		Token { kind: TokenKindRBrace },
		Token { kind: TokenKindNewline },
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
		Token { kind: TokenKindNewline },
	}, test)
}

func TestTokenizeNumbers (test *testing.T) {
	checkTokenSlice("../tests/lexer/numbers.arf", []Token {
		Token { kind: TokenKindUInt, value: uint64(0) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindUInt, value: uint64(8) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindUInt, value: uint64(83628266) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindUInt, value: uint64(83628266) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindUInt, value: uint64(83628266) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindUInt, value: uint64(83628266) },
		Token { kind: TokenKindNewline },
		
		Token { kind: TokenKindInt, value: int64(-83628266) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindInt, value: int64(-83628266) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindInt, value: int64(-83628266) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindInt, value: int64(-83628266) },
		Token { kind: TokenKindNewline },
		
		Token { kind: TokenKindFloat, value: float64(0.123478) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindFloat, value: float64(234.3095) },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindFloat, value: float64(-2.312) },
		Token { kind: TokenKindNewline },
	}, test)
}

func TestTokenizeText (test *testing.T) {
	checkTokenSlice("../tests/lexer/text.arf", []Token {
		Token { kind: TokenKindString, value: "hello world!\a\b\f\n\r\t\v'\"\\" },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindRune, value: '\a' },
		Token { kind: TokenKindRune, value: '\b' },
		Token { kind: TokenKindRune, value: '\f' },
		Token { kind: TokenKindRune, value: '\n' },
		Token { kind: TokenKindRune, value: '\r' },
		Token { kind: TokenKindRune, value: '\t' },
		Token { kind: TokenKindRune, value: '\v' },
		Token { kind: TokenKindRune, value: '\'' },
		Token { kind: TokenKindRune, value: '"'  },
		Token { kind: TokenKindRune, value: '\\' },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindString, value: "hello world \x40\u0040\U00000040!" },
		Token { kind: TokenKindNewline },
	}, test)
}

func TestTokenizeIndent (test *testing.T) {
	checkTokenSlice("../tests/lexer/indent.arf", []Token {
		Token { kind: TokenKindName, value: "line1" },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindIndent, value: 1 },
		Token { kind: TokenKindName, value: "line2" },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindIndent, value: 4 },
		Token { kind: TokenKindName, value: "line3" },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindName, value: "line4" },
		Token { kind: TokenKindNewline },
		Token { kind: TokenKindIndent, value: 2 },
		Token { kind: TokenKindName, value: "line5" },
		Token { kind: TokenKindNewline },
	}, test)
}
