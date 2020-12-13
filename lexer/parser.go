package lexer

import "github.com/deckarep/golang-set"

type precedence struct {
	value     int
	leftAssoc bool
}

func newPrecedence(value int, leftAssoc bool) *precedence {
	return &precedence{
		value, leftAssoc,
	}
}

const (
	LEFT  bool = true
	RIGHT bool = false
)

// Operators 操作
type Operators struct {
	opMap map[string]*precedence
}

// NewOperators 创建Operators对象
func NewOperators() Operators {
	return Operators{make(map[string]*precedence)}
}

// Add 添加操作
func (o Operators) Add(name string, prec int, leftAssoc bool) {
	o.opMap[name] = newPrecedence(prec, leftAssoc)
}

// BasicParser 语法解析器
type BasicParser struct {
	reserved  mapset.Set
	operators Operators
	parser    Parser
	primary   Parser
}

// NewBasicParser 创建Parser对象
func NewBasicParser() BasicParser {
	reserved := mapset.NewSet(";", "}", EOL)
	return BasicParser{
		reserved:  reserved,
		operators: NewOperators(),
		parser:    Rule(),
	}
}
