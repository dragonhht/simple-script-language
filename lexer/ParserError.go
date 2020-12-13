package lexer

import "fmt"

// ParserError 解析错误
func ParserError(msg string, token Token) {
	panic(fmt.Sprintf("syntax error around %v. %v", errorLocation(token), msg))
}

// errorLocation 错误定位信息
func errorLocation(token Token) string {
	if token == EOF {
		return "the last line"
	} else {
		return fmt.Sprintf(`"%v" at line %v`, token.GetText(), token.GetLineNumber())
	}
}
