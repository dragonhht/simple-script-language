package lexer

import (
	"github.com/deckarep/golang-set"
	"simple-script-language/utils/list"
)

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
	primary := RuleByType(NewPrimaryExpr(list.New(0))).Or([]*Parser{
		Rule().Sep("(").Ast(expr0).Sep(")"),
		Rule().Number(NewNumberNode(nil)),
		Rule().Identifier(NewVariableNode(nil), reserved),
		Rule().String(NewStringNode(nil)),
	})
	factor := Rule().Or([]*Parser{
		RuleByType(NewNegativeExprNode(list.New(0))).Sep("-").Ast(primary),
		primary,
	})
	expr := expr0.Expression(NewBinaryExprNode(list.New(0)), factor, operators)
	statement0 := Rule()
	block := RuleByType(NewBlockStatementNode(list.New(0))).Sep("{").Option(statement0).Repeat(Rule().Sep(";", EOL).Option(statement0)).Sep("}")
	simple := RuleByType(NewPrimaryExpr(list.New(0))).Ast(expr)
	statement := statement0.Or([]*Parser{
		RuleByType(NewIfStatementNode(list.New(0))).Sep("if").Ast(expr).Ast(block).Option(
			Rule().Sep("else").Ast(block)),
		RuleByType(NewWhileStatementNode(list.New(0))).Sep("while").Ast(expr).Ast(block),
		simple,
	})
	program := Rule().Or([]*Parser{
		statement,
		RuleByType(NewNullStatementNode(list.New(0))),
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

// FuncParser 函数解析器
type FuncParser struct {
	BasicParser
	param     *Parser
	params    *Parser
	paramList *Parser
	def       *Parser
	args      *Parser
	postfix   *Parser
}

// NewFuncParser 创建FuncParser
func NewFuncParser() FuncParser {
	bp := NewBasicParser()
	param := Rule().Identifier(nil, bp.reserved)
	params := RuleByType(NewParameterListNode(list.New(0))).Ast(param).Repeat(Rule().Sep(",").Ast(param))
	paramList := Rule().Sep("(").Maybe(params).Sep(")")
	def := RuleByType(NewDefStatementNode(list.New(0))).Sep("def").Identifier(nil, bp.reserved).Ast(paramList).Ast(bp.block)
	args := RuleByType(NewArgumentsNode(list.New(0))).Ast(bp.expr).Repeat(Rule().Sep(",").Ast(bp.expr))
	postfix := Rule().Sep("(").Maybe(args).Sep(")")

	bp.reserved.Add(")")
	bp.primary.Repeat(postfix)
	bp.simple.Option(args)
	bp.program.InsertChoice(def)
	return FuncParser{
		BasicParser: bp,
		param:       param,
		params:      params,
		paramList:   paramList,
		def:         def,
		args:        args,
		postfix:     postfix,
	}
}
