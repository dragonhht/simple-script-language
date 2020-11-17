package lexer

import (
	"bufio"
	"container/list"
	"regexp"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||\p{Punct})?` // 匹配的正则表达式

// Lexer 词法分析器
type Lexer struct {
	pattern *regexp.Regexp // 正则对象
	queue *list.List // 单词暂存列表
	hasMore bool // 是否还有为解析单词
	reader *bufio.Reader // 内容读取器
}

// NewLexer 创建Lexer对象
func NewLexer(reader *bufio.Reader) *Lexer {
	pattern, _ := regexp.Compile(regexPat)
    return &Lexer{
		pattern: pattern,
		queue:   list.New(),
		hasMore: true,
		reader:  reader,
	}
}