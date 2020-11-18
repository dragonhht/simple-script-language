package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||\p{Punct})?` // 匹配的正则表达式

// Lexer 词法分析器
type Lexer struct {
	pattern *regexp.Regexp // 正则对象
	queue []Token // 单词暂存列表
	hasMore bool // 是否还有为解析单词
	reader *bufio.Scanner // 内容读取器
}

// NewLexer 创建Lexer对象
func NewLexer(reader *bufio.Scanner) *Lexer {
	pattern, _ := regexp.Compile(regexPat)
    return &Lexer{
		pattern: pattern,
		queue:   make([]Token, 10),
		hasMore: true,
		reader:  reader,
	}
}

// ParseError 解析异常
func ParseError() error {
    return errors.New("")
}

// Read 从源代码源头逐一获取单词
func (l *Lexer) Read() (Token, error) {
	fill, err := l.fillQueue(0)
	if err != nil {
		return nil, err
	}
	if fill {
		token := l.queue[0]
		l.queue = l.queue[1:]
		return token, nil
	} else {
		return EOF, nil
	}
}

// fillQueue 将单词加载到暂存列表
func (l *Lexer) fillQueue(index int) (bool, error) {
	// 当需要获取超过缓存长度的单词时，进行加载
	for {
		if index < len(l.queue) {
			break
		}
		if l.hasMore {
			err := l.readLine()
			if err != nil {
				return false, err
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}

// readLine 逐行读取单词
func (l *Lexer) readLine() error {
	var line string
	if l.reader.Scan() {
		line = l.reader.Text()
	} else {
		// 已到最后，没有更多
		l.hasMore = false
		return nil
	}
	// TODO 如何获取行号
	var lineNo = 0
	pos := 0
	endPos := len(line)
	for {
		matcherLine := line[pos:]
		if pos >= endPos {
			break
		}
		loc := l.pattern.FindIndex([]byte(matcherLine))
		// 起始匹配
		if loc[0] == 0 {
			l.addToken(lineNo, matcherLine)
			pos = loc[1]
		} else {
			return errors.New(fmt.Sprintf("bad token at line %d", lineNo))
		}
	}
	l.queue = append(l.queue, NewIdToken(lineNo, EOL))
	return nil
}

// addToken 创建并保存Token对象
func (l *Lexer) addToken(lineNo int, lineStr string) {

}