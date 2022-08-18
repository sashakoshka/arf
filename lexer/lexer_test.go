package lexer

import "testing"
import "git.tebibyte.media/sashakoshka/arf/file"
import "git.tebibyte.media/sashakoshka/arf/types"
import "git.tebibyte.media/sashakoshka/arf/infoerr"

func quickToken (kind TokenKind, value any) (token Token) {
	return Token { kind: kind, value: value }
}

func checkTokenSlice (filePath string, test *testing.T, correct ...Token) {
	test.Log("checking lexer results for", filePath)
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

func compareErr (
	filePath string,
	correctKind    infoerr.ErrorKind,
	correctMessage string,
	correctRow     int,
	correctColumn  int,
	correctWidth   int,
	test *testing.T,
) {
	test.Log("testing errors in", filePath)
	file, err := file.Open(filePath)
	if err != nil {
		test.Log(err)
		test.Fail()
		return
	}
	
	_, err = Tokenize(file)
	check := err.(infoerr.Error)

	test.Log("error that was recieved:")
	test.Log(check)

	if check.Kind() != correctKind {
		test.Log("mismatched error kind")
		test.Log("- want:", correctKind)
		test.Log("- have:", check.Kind())
		test.Fail()
	}

	if check.Message() != correctMessage {
		test.Log("mismatched error message")
		test.Log("- want:", correctMessage)
		test.Log("- have:", check.Message())
		test.Fail()
	}

	if check.Row() != correctRow {
		test.Log("mismatched error row")
		test.Log("- want:", correctRow)
		test.Log("- have:", check.Row())
		test.Fail()
	}

	if check.Column() != correctColumn {
		test.Log("mismatched error column")
		test.Log("- want:", correctColumn)
		test.Log("- have:", check.Column())
		test.Fail()
	}

	if check.Width() != correctWidth {
		test.Log("mismatched error width")
		test.Log("- want:", check.Width())
		test.Log("- have:", correctWidth)
		test.Fail()
	}
}

func TestTokenizeAll (test *testing.T) {
	checkTokenSlice("../tests/lexer/all.arf", test,
		quickToken(TokenKindSeparator, nil),
		quickToken (TokenKindPermission, types.Permission {
			Internal: types.ModeRead,
			External: types.ModeWrite,
		}),
		quickToken(TokenKindReturnDirection, nil),
		quickToken(TokenKindInt, int64(-349820394)),
		quickToken(TokenKindUInt, uint64(932748397)),
		quickToken(TokenKindFloat, 239485.37520),
		quickToken(TokenKindString, "hello world!\n"),
		quickToken(TokenKindRune, 'E'),
		quickToken(TokenKindName, "helloWorld"),
		quickToken(TokenKindColon, nil),
		quickToken(TokenKindDot, nil),
		quickToken(TokenKindComma, nil),
		quickToken(TokenKindElipsis, nil),
		quickToken(TokenKindLBracket, nil),
		quickToken(TokenKindRBracket, nil),
		quickToken(TokenKindLBrace, nil),
		quickToken(TokenKindRBrace, nil),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindPlus, nil),
		quickToken(TokenKindMinus, nil),
		quickToken(TokenKindIncrement, nil),
		quickToken(TokenKindDecrement, nil),
		quickToken(TokenKindAsterisk, nil),
		quickToken(TokenKindSlash, nil),
		quickToken(TokenKindAt, nil),
		quickToken(TokenKindExclamation, nil),
		quickToken(TokenKindPercent, nil),
		quickToken(TokenKindTilde, nil),
		quickToken(TokenKindLessThan, nil),
		quickToken(TokenKindLShift, nil),
		quickToken(TokenKindGreaterThan, nil),
		quickToken(TokenKindRShift, nil),
		quickToken(TokenKindBinaryOr, nil),
		quickToken(TokenKindLogicalOr, nil),
		quickToken(TokenKindBinaryAnd, nil),
		quickToken(TokenKindLogicalAnd, nil),
		quickToken(TokenKindNewline, nil),
	)
}

func TestTokenizeNumbers (test *testing.T) {
	checkTokenSlice("../tests/lexer/numbers.arf", test,
		quickToken(TokenKindUInt, uint64(0)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindUInt, uint64(8)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindUInt, uint64(83628266)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindUInt, uint64(83628266)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindUInt, uint64(83628266)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindUInt, uint64(83628266)),
		quickToken(TokenKindNewline, nil),
		
		quickToken(TokenKindInt, int64(-83628266)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindInt, int64(-83628266)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindInt, int64(-83628266)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindInt, int64(-83628266)),
		quickToken(TokenKindNewline, nil),
		
		quickToken(TokenKindFloat, float64(0.123478)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindFloat, float64(234.3095)),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindFloat, float64(-2.312)),
		quickToken(TokenKindNewline, nil),
	)
}

func TestTokenizeText (test *testing.T) {
	checkTokenSlice("../tests/lexer/text.arf", test,
		quickToken(TokenKindString, "hello world!\a\b\f\n\r\t\v'\"\\"),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindRune, '\a'),
		quickToken(TokenKindRune, '\b'),
		quickToken(TokenKindRune, '\f'),
		quickToken(TokenKindRune, '\n'),
		quickToken(TokenKindRune, '\r'),
		quickToken(TokenKindRune, '\t'),
		quickToken(TokenKindRune, '\v'),
		quickToken(TokenKindRune, '\''),
		quickToken(TokenKindRune, '"' ),
		quickToken(TokenKindRune, '\\'),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindString, "hello world \x40\u0040\U00000040!"),
		quickToken(TokenKindNewline, nil),
	)
}

func TestTokenizeIndent (test *testing.T) {
	checkTokenSlice("../tests/lexer/indent.arf", test,
		quickToken(TokenKindName, "line1"),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindIndent, 1),
		quickToken(TokenKindName, "line2"),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindIndent, 4),
		quickToken(TokenKindName, "line3"),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindName, "line4"),
		quickToken(TokenKindNewline, nil),
		quickToken(TokenKindIndent, 2),
		quickToken(TokenKindName, "line5"),
		quickToken(TokenKindNewline, nil),
	)
}

func TestTokenizeErr (test *testing.T) {
	compareErr (
		"../tests/lexer/error/unexpectedSymbol.arf",
		infoerr.ErrorKindError,
		"unexpected symbol character ;",
		1, 5, 1,
		test)
	
	compareErr (
		"../tests/lexer/error/excessDataRune.arf",
		infoerr.ErrorKindError,
		"excess data in rune literal",
		1, 1, 7,
		test)
	
	compareErr (
		"../tests/lexer/error/unknownEscape.arf",
		infoerr.ErrorKindError,
		"unknown escape character g",
		1, 2, 1,
		test)
}
