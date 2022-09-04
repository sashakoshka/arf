package lexer

import "testing"
import "git.tebibyte.media/arf/arf/file"
import "git.tebibyte.media/arf/arf/types"
import "git.tebibyte.media/arf/arf/infoerr"

func quickToken (width int, kind TokenKind, value any) (token Token) {
	token.location.SetWidth(width)
	token.kind  = kind
	token.value = value
	return
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
		if token.location.Width() != correct[index].location.Width() {
			test.Log("token", index, "has bad width")
			test.Log (
				"have", token.location.Width(),
				"want", correct[index].location.Width())
			test.Fail()
			return
		}
		
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
		quickToken(3, TokenKindSeparator, nil),
		quickToken(2, TokenKindPermission, types.PermissionReadWrite),
		quickToken(2, TokenKindReturnDirection, nil),
		quickToken(10, TokenKindInt, int64(-349820394)),
		quickToken(9, TokenKindUInt, uint64(932748397)),
		quickToken(12, TokenKindFloat, 239485.37520),
		quickToken(16, TokenKindString, "hello world!\n"),
		quickToken(3, TokenKindRune, 'E'),
		quickToken(10, TokenKindName, "helloWorld"),
		quickToken(1, TokenKindColon, nil),
		quickToken(1, TokenKindDot, nil),
		quickToken(1, TokenKindComma, nil),
		quickToken(2, TokenKindElipsis, nil),
		quickToken(1, TokenKindLBracket, nil),
		quickToken(1, TokenKindRBracket, nil),
		quickToken(1, TokenKindLBrace, nil),
		quickToken(1, TokenKindRBrace, nil),
		quickToken(1, TokenKindNewline, nil),
		quickToken(1, TokenKindPlus, nil),
		quickToken(1, TokenKindMinus, nil),
		quickToken(2, TokenKindIncrement, nil),
		quickToken(2, TokenKindDecrement, nil),
		quickToken(1, TokenKindAsterisk, nil),
		quickToken(1, TokenKindSlash, nil),
		quickToken(1, TokenKindAt, nil),
		quickToken(1, TokenKindExclamation, nil),
		quickToken(1, TokenKindPercent, nil),
		quickToken(2, TokenKindPercentAssignment, nil),
		quickToken(1, TokenKindTilde, nil),
		quickToken(2, TokenKindTildeAssignment, nil),
		quickToken(1, TokenKindAssignment, nil),
		quickToken(2, TokenKindEqualTo, nil),
		quickToken(2, TokenKindNotEqualTo, nil),
		quickToken(1, TokenKindLessThan, nil),
		quickToken(2, TokenKindLessThanEqualTo, nil),
		quickToken(2, TokenKindLShift, nil),
		quickToken(3, TokenKindLShiftAssignment, nil),
		quickToken(1, TokenKindGreaterThan, nil),
		quickToken(2, TokenKindGreaterThanEqualTo, nil),
		quickToken(2, TokenKindRShift, nil),
		quickToken(3, TokenKindRShiftAssignment, nil),
		quickToken(1, TokenKindBinaryOr, nil),
		quickToken(2, TokenKindBinaryOrAssignment, nil),
		quickToken(2, TokenKindLogicalOr, nil),
		quickToken(1, TokenKindBinaryAnd, nil),
		quickToken(2, TokenKindBinaryAndAssignment, nil),
		quickToken(2, TokenKindLogicalAnd, nil),
		quickToken(1, TokenKindBinaryXor, nil),
		quickToken(2, TokenKindBinaryXorAssignment, nil),
		quickToken(1, TokenKindNewline, nil),
	)
}

func TestTokenizeNumbers (test *testing.T) {
	checkTokenSlice("../tests/lexer/numbers.arf", test,
		quickToken(1, TokenKindUInt, uint64(0)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(1, TokenKindUInt, uint64(8)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(8, TokenKindUInt, uint64(83628266)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(29, TokenKindUInt, uint64(83628266)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(9, TokenKindUInt, uint64(83628266)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(10, TokenKindUInt, uint64(83628266)),
		quickToken(1, TokenKindNewline, nil),
		
		quickToken(9, TokenKindInt, int64(-83628266)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(30, TokenKindInt, int64(-83628266)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(10, TokenKindInt, int64(-83628266)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(11, TokenKindInt, int64(-83628266)),
		quickToken(1, TokenKindNewline, nil),
		
		quickToken(8, TokenKindFloat, float64(0.123478)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(8, TokenKindFloat, float64(234.3095)),
		quickToken(1, TokenKindNewline, nil),
		quickToken(6, TokenKindFloat, float64(-2.312)),
		quickToken(1, TokenKindNewline, nil),
	)
}

func TestTokenizeText (test *testing.T) {
	checkTokenSlice("../tests/lexer/text.arf", test,
		quickToken(34, TokenKindString, "hello world!\a\b\f\n\r\t\v'\"\\"),
		quickToken(1, TokenKindNewline, nil),
		quickToken(4, TokenKindRune, '\a'),
		quickToken(4, TokenKindRune, '\b'),
		quickToken(4, TokenKindRune, '\f'),
		quickToken(4, TokenKindRune, '\n'),
		quickToken(4, TokenKindRune, '\r'),
		quickToken(4, TokenKindRune, '\t'),
		quickToken(4, TokenKindRune, '\v'),
		quickToken(4, TokenKindRune, '\''),
		quickToken(4, TokenKindRune, '"' ),
		quickToken(4, TokenKindRune, '\\'),
		quickToken(1, TokenKindNewline, nil),
		quickToken(35, TokenKindString, "hello world \x40\u0040\U00000040!"),
		quickToken(1, TokenKindNewline, nil),
	)
}

func TestTokenizeIndent (test *testing.T) {
	checkTokenSlice("../tests/lexer/indent.arf", test,
		quickToken(5, TokenKindName, "line1"),
		quickToken(1, TokenKindNewline, nil),
		quickToken(1, TokenKindIndent, 1),
		quickToken(5, TokenKindName, "line2"),
		quickToken(1, TokenKindNewline, nil),
		quickToken(4, TokenKindIndent, 4),
		quickToken(5, TokenKindName, "line3"),
		quickToken(1, TokenKindNewline, nil),
		quickToken(5, TokenKindName, "line4"),
		quickToken(1, TokenKindNewline, nil),
		quickToken(2, TokenKindIndent, 2),
		quickToken(5, TokenKindName, "line5"),
		quickToken(1, TokenKindNewline, nil),
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
