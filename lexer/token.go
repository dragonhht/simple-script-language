// lexer 词法分析器
package lexer

import "errors"

var (
	EOF = NewToken(1)
	EOL = "\\n"
)

// Token 词法分析的结果(单词)
type Token struct {
	lineNumber int //行号
}

// NewToken 创建一个新的Token
func NewToken(line int) *Token {
	return &Token{lineNumber: line}
}

// GetLineNumber 获取行号
func (t *Token) GetLineNumber() int {
    return t.lineNumber
}

// IsIdentifier 是否为标识符(变量名、函数名、类名)
func (t *Token) IsIdentifier() bool {
    return false
}

// IsNumber 是否为整型字面量
func (t *Token) IsNumber() bool {
    return false
}

// IsString 是否为字符串字面量
func (t *Token) IsString() bool {
    return false
}

// GetNumber 获取整型字面量的值
func (t *Token) GetNumber() (int, error) {
	return -1, errors.New("not number token")
}

// GetText 获取字符串字面量的值
func (t *Token) GetText() string {
	return ""
}
