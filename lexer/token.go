// lexer 词法分析器
package lexer

import (
	"errors"
	"strconv"
)

// TokenInterface 单词接口
type Token interface {
	GetLineNumber() int      // 获取行号
	IsIdentifier() bool      // 是否为标识符(变量名、函数名、类名)
	IsNumber() bool          // 是否为整型字面量
	IsString() bool          // 是否为字符串字面量
	GetNumber() (int, error) // 获取整型字面量的值
	GetText() string         // 获取字符串字面量的值
}

var (
	EOF = NewToken(1)
	EOL = "\\n"
)

// Token 词法分析的结果(单词)
type AbstractToken struct {
	lineNumber int //行号
}

// NewToken 创建一个新的Token
func NewToken(line int) AbstractToken {
	return AbstractToken{lineNumber: line}
}

// GetLineNumber 获取行号
func (t AbstractToken) GetLineNumber() int {
	return t.lineNumber
}

// IsIdentifier 是否为标识符(变量名、函数名、类名)
func (t AbstractToken) IsIdentifier() bool {
	return false
}

// IsNumber 是否为整型字面量
func (t AbstractToken) IsNumber() bool {
	return false
}

// IsString 是否为字符串字面量
func (t AbstractToken) IsString() bool {
	return false
}

// GetNumber 获取整型字面量的值
func (t AbstractToken) GetNumber() (int, error) {
	return -1, errors.New("not number token")
}

// GetText 获取字符串字面量的值
func (t AbstractToken) GetText() string {
	return ""
}

// NumToken 整型字面量的Token
type NumToken struct {
	AbstractToken
	value int // 整型字面量的值
}

// NewNumToken 创建NumToken对象
func NewNumToken(line, value int) NumToken {
	return NumToken{
		AbstractToken: NewToken(line),
		value:         value,
	}
}

// IsNumber 是否为整型字面量
func (n NumToken) IsNumber() bool {
	return true
}

// GetText 获取字符串字面量的值
func (n NumToken) GetText() string {
	return strconv.Itoa(n.value)
}

// GetNumber 获取整型字面量的值
func (n NumToken) GetNumber() (int, error) {
	return n.value, nil
}

// IdToken 标志符类型的Token
type IdToken struct {
	AbstractToken
	text string // 标识
}

// NewIdToken 创建IdToken对象
func NewIdToken(line int, id string) IdToken {
	return IdToken{
		AbstractToken: NewToken(line),
		text:          id,
	}
}

// IsIdentifier 是否为标识符(变量名、函数名、类名)
func (i IdToken) IsIdentifier() bool {
	return true
}

// GetText 获取字符串字面量的值
func (i IdToken) GetText() string {
	return i.text
}

// StrToken 字符串字面量的Token
type StrToken struct {
	AbstractToken
	literal string // 字符串值
}

// NewStrToken 创建StrToken对象
func NewStrToken(line int, literal string) StrToken {
	return StrToken{
		AbstractToken: NewToken(line),
		literal:       literal,
	}
}

// IsString 是否为字符串字面量
func (s StrToken) IsString() bool {
	return true
}

// GetText 获取字符串字面量的值
func (s StrToken) GetText() string {
	return s.literal
}
