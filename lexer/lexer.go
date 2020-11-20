package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// regexPat 正则表达式,使用命名分组
const regexPat = `\s*((?P<notes>//.*)|(?P<number>[0-9]+)|(?P<stringVal>"(\\"|\\\\|\\n|[^"])*")|(?P<string>[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]]))?` // 匹配的正则表达式

// Lexer 词法分析器
type Lexer struct {
	pattern *regexp.Regexp // 正则对象
	queue   []Token        // 单词暂存列表
	hasMore bool           // 是否还有为解析单词
	reader  *bufio.Scanner // 内容读取器
	lineNo  int            // 行号
}

// NewLexer 创建Lexer对象
func NewLexer(reader *bufio.Scanner) *Lexer {
	pattern, _ := regexp.Compile(regexPat)
	return &Lexer{
		pattern: pattern,
		queue:   make([]Token, 0),
		hasMore: true,
		reader:  reader,
	}
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

// Peek 读取Read读取的单词n位后的单词
func (l *Lexer) Peek(n int) (Token, error) {
	fill, err := l.fillQueue(n)
	if err != nil {
		return nil, err
	}
	if fill {
		return l.queue[n], nil
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
		l.lineNo++
		line = l.reader.Text()
	} else {
		// 已到最后，没有更多
		l.hasMore = false
		return nil
	}
	pos := 0
	endPos := len(line)
	matcherLine := line
	for {
		matcherLine := matcherLine[pos:]
		if pos >= endPos {
			break
		}
		loc := l.pattern.FindIndex([]byte(matcherLine))
		// 起始匹配
		if loc[0] == 0 {
			l.addToken(l.lineNo, matcherLine)
			pos += loc[1]
		} else {
			return errors.New(fmt.Sprintf("bad token at line %d", l.lineNo))
		}
	}
	l.queue = append(l.queue, NewIdToken(l.lineNo, EOL))
	return nil
}

// addToken 创建并保存Token对象
func (l *Lexer) addToken(lineNo int, lineStr string) {
	// 命名分组
	match := l.pattern.FindStringSubmatch(lineStr)
	// groupNames := l.pattern.SubexpNames()
	m := match[1]
	if m != "" {
		if match[2] == "" {
			var token Token
			if match[3] != "" {
				value, _ := strconv.Atoi(m)
				token = NewNumToken(lineNo, value)
			} else if match[4] != "" {
				token = NewStrToken(lineNo, toStringLiteral(m))
			} else {
				token = NewIdToken(lineNo, m)
			}
			l.queue = append(l.queue, token)
		}
	}
}

// toStringLiteral 字符串类型的Token转换字符格式
func toStringLiteral(str string) string {
	var buf strings.Builder
	len := len(str) - 1
	for i := 1; i < len; i++ {
		c := str[i]
		if c == '\\' && i+1 < len {
			c2 := str[i+1]
			if c2 == '"' || c2 == '\\' {
				i += 1
				c = str[i]
			} else if c2 == 'n' {
				i += 1
				c = '\n'
			}
		}
		buf.WriteByte(c)
	}
	return buf.String()
}
