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
	reserved   mapset.Set
	operators  Operators
	parser     *Parser
	primary    *Parser
	factor     *Parser
	expr       *Parser
	statement0 *Parser
	block      *Parser
	simple     *Parser
	statement  *Parser
	program    *Parser
}

// NewBasicParser 创建Parser对象
func NewBasicParser() BasicParser {
	reserved := mapset.NewSet(";", "}", EOL)
	operators := NewOperators()
	operators.Add("=", 1, RIGHT)
	operators.Add("==", 2, LEFT)
	operators.Add(">", 2, LEFT)
	operators.Add("<", 2, LEFT)
	operators.Add("+", 3, LEFT)
	operators.Add("-", 3, LEFT)
	operators.Add("*", 4, LEFT)
	operators.Add("/", 4, LEFT)
	operators.Add("%", 4, LEFT)

	expr0 := Rule()
	primary := RuleByType(NewPrimaryExpr([]TreeNode{})).Or([]*Parser{
		Rule().Sep("(").Ast(expr0).Sep(")"),
		Rule().Number(NewNumberNode(nil)),
		Rule().Identifier(NewVariableNode(nil), reserved),
		Rule().String(NewStringNode(nil)),
	})
	factor := Rule().Or([]*Parser{
		RuleByType(NewNegativeExprNode([]TreeNode{})).Sep("-").Ast(primary),
		primary,
	})
	expr := expr0.Expression(NewBinaryExprNode([]TreeNode{}), factor, operators)
	statement0 := Rule()
	block := RuleByType(NewBlockStatementNode([]TreeNode{})).Sep("{").Option(statement0).Repeat(Rule().Sep(";", EOL).Option(statement0)).Sep("}")
	simple := RuleByType(NewPrimaryExpr([]TreeNode{})).Ast(expr)
	statement := statement0.Or([]*Parser{
		RuleByType(NewIfStatementNode([]TreeNode{})).Sep("if").Ast(expr).Ast(block).Option(
			Rule().Sep("else").Ast(block)),
		RuleByType(NewWhileStatementNode([]TreeNode{})).Sep("while").Ast(expr).Ast(block),
		simple,
	})
	program := Rule().Or([]*Parser{
		statement,
		RuleByType(NewNullStatementNode([]TreeNode{})),
	}).Sep(";", EOL)
	return BasicParser{
		reserved:   reserved,
		operators:  operators,
		parser:     expr0,
		primary:    primary,
		factor:     factor,
		expr:       expr,
		statement0: statement0,
		block:      block,
		simple:     simple,
		statement:  statement,
		program:    program,
	}
}

// Parser 解析
func (b BasicParser) Parser(lexer *Lexer) TreeNode {
	return b.program.parse(lexer)
}
