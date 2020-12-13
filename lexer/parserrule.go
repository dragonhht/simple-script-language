package lexer

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
)

// ParserElement 解析器元素接口
type ParserElement interface {
	Parse(lexer Lexer, res []TreeNode) []TreeNode // 解析
	Match(lexer Lexer) bool                       // 匹配
}

// Leaf 叶子节点解析器元素
type Leaf struct {
	tokens []string // 单词
}

// NewLeaf 创建Leaf对象
func NewLeaf(pat []string) Leaf {
	return Leaf{pat}
}

// Parse 解析
func (l Leaf) Parse(lexer Lexer, res []TreeNode) []TreeNode {
	t, err := lexer.Read()
	if err != nil {
		panic(err)
	}
	if t.IsIdentifier() {
		for _, v := range l.tokens {
			if v == t.GetText() {
				return l.find(res, t)
			}
		}
	}
	if len(l.tokens) > 0 {
		ParserError(l.tokens[0]+" expected.", t)
	} else {
		ParserError("", t)
	}
	return nil
}

// find 查找单词元素
func (l Leaf) find(res []TreeNode, token Token) []TreeNode {
	return append(res, NewLeafNode(token))
}

// Match 匹配
func (l Leaf) Match(lexer Lexer) bool {
	t, err := lexer.Peek(0)
	if err != nil {
		panic(err)
	}
	if t.IsIdentifier() {
		for _, v := range l.tokens {
			if v == t.GetText() {
				return true
			}
		}
	}
	return false
}

// Skip 跳过元素
type Skip struct {
	Leaf
}

// NewSkip 创建Skip对象
func NewSkip(pat []string) Skip {
	return Skip{NewLeaf(pat)}
}

// find 查找
func find(res []TreeNode, token Token) []TreeNode {
	return nil
}

// TreeParser 树解析器
type TreeParser struct {
	parser Parser
}

// NewTreeParser 创建TreeParser对象
func NewTreeParser(p Parser) TreeParser {
	return TreeParser{p}
}

// Parse 解析
func (t TreeParser) Parse(lexer Lexer, res []TreeNode) []TreeNode {
	return append(res, t.parser.parse(lexer))
}

// Match 匹配
func (t TreeParser) Match(lexer Lexer) bool {
	return t.parser.Match(lexer)
}

// AToken
type AToken struct {
	factory *Factory
	test    func(token Token) bool
}

// NewAToken
func NewAToken(node interface{}) AToken {
	if node == nil {
		node = NewLeafNode(nil)
	}
	factory := getFactory(node)
	return AToken{factory: factory}
}

// Parse 解析
func (a AToken) Parse(lexer Lexer, res []TreeNode) []TreeNode {
	t, err := lexer.Read()
	if err != nil {
		panic(err)
	}
	if a.test(t) {
		leaf := a.factory.make(t)
		return append(res, leaf)
	} else {
		ParserError("", t)
		return nil
	}
}

// Match 匹配
func (a AToken) Match(lexer Lexer) bool {
	t, err := lexer.Peek(0)
	if err != nil {
		panic(err)
	}
	return a.test(t)
}

// IdTokenParser 标识解析器
type IdTokenParser struct {
	AToken
	reserved mapset.Set
}

// NewIdTokenParser 创建IdTokenParser对象
func NewIdTokenParser(typ interface{}, r mapset.Set) IdTokenParser {
	aToken := NewAToken(typ)
	reserved := r
	if r == nil {
		reserved = mapset.NewSet()
	}
	idToken := IdTokenParser{
		AToken:   aToken,
		reserved: reserved,
	}
	idToken.test = func(token Token) bool {
		return token.IsIdentifier() && !idToken.reserved.Contains(token.GetText())
	}
	return idToken
}

// NumTokenParser 数值解析器
type NumTokenParser struct {
    AToken
}

// NewNumTokenParser 创建NumTokenParser
func NewNumTokenParser(typ interface{}) NumTokenParser {
	aToken := NewAToken(typ)
	numToken := NumTokenParser{aToken}
	numToken.test = func(token Token) bool {
		return token.IsNumber()
	}
	return numToken
}

// StrTokenParser 字符解析器
type StrTokenParser struct {
    AToken
}

// NewStrTokenParser 创建StrTokenParser
func NewStrTokenParser(typ interface{}) StrTokenParser {
    aToken := NewAToken(typ)
    strToken := StrTokenParser{aToken}
    strToken.test = func(token Token) bool {
		return token.IsString()
	}
    return strToken
}

// ExprParser 表达式元素解析器
type ExprParser struct {
    factory *Factory
    ops Operators
    parser Parser
}

// NewExprParser 创建ExprParser
func NewExprParser(typ interface{}, exp Parser, ops Operators) ExprParser {
	return ExprParser{
		factory: getForASTListFactory(typ),
		ops:     ops,
		parser:  exp,
	}
}

// Parse 解析
func (e ExprParser) Parse(lexer Lexer, res []TreeNode) []TreeNode {
	right := e.parser.parse(lexer)
	prec := e.nextOperator(lexer)
	for prec != nil {
		right = e.doShift(lexer, right, prec.value)
		prec = e.nextOperator(lexer)
	}
	return append(res, right)
}

func (e ExprParser) doShift(lexer Lexer, left TreeNode, prec int) TreeNode {
	tree := make([]TreeNode, 2)
	tree[0] = left
	t, err := lexer.Read()
	if err != nil {
		panic(err)
	}
	tree[1] = NewLeafNode(t)
	right := e.parser.parse(lexer);
	next := e.nextOperator(lexer)
	for next != nil && e.rightIsExpr(prec, next) {
		right = e.doShift(lexer, right, next.value)
	}
	tree = append(tree, right)
	return e.factory.make(tree)
}

// nextOperator 下一个操作
func (e ExprParser) nextOperator(lexer Lexer) *precedence {
	t, err := lexer.Read()
	if err != nil {
		panic(err)
	}
	if t.IsIdentifier() {
		return e.ops.opMap[t.GetText()]
	} else {
		return nil
	}
}

// rightIsExpr 判断右边是否为表达式
func (e ExprParser) rightIsExpr(prec int, next *precedence) bool {
	if next.leftAssoc {
		return prec < next.value
	} else  {
		return prec <= next.value
	}
}

// Match 匹配
func (e ExprParser) Match(lexer Lexer) bool {
	return e.parser.Match(lexer)
}

// OrTree Or逻辑解析器元素
type OrTree struct {
	parsers []Parser
}

// NewOrTree 创建Or逻辑解析器元素
func NewOrTree(parsers []Parser) OrTree {
	return OrTree{parsers: parsers}
}

// Parse Or逻辑解析
func (o OrTree) Parse(lexer Lexer, res []TreeNode) []TreeNode {
	p := o.choose(lexer)
	if &p == nil {
		panic(fmt.Sprintf("解析错误: %v\n", lexer.Peek(0)))
	} else {
		res = append(res, p.parse(lexer))
	}
	return res
}

// Match or逻辑匹配
func (o OrTree) Match(lexer Lexer) bool {
	p := o.choose(lexer)
	return p != nil
}

// choose
func (o OrTree) choose(lexer Lexer) *Parser {
	for _, p := range o.parsers {
		if p.Match(lexer) {
			return &p
		}
	}
	return nil
}

// insert
func (o OrTree) insert(p Parser) {
	ps := make([]Parser, 1)
	ps[0] = p
	o.parsers = append(ps, o.parsers...)
}

// RepeatParser
type RepeatParser struct {
    parser Parser
    onlyOnce bool
}

// NewRepeatParser
func NewRepeatParser(parser Parser, onlyOnce bool) RepeatParser {
    return RepeatParser {
    	parser,
    	onlyOnce,
	}
}

// Parse 解析
func (r RepeatParser) Parse(lexer Lexer, res []TreeNode) []TreeNode {
	for r.parser.Match(lexer) {
		t := r.parser.parse(lexer)
		switch t.(type) {
		case BranchNode:
			break
		default:
			if t.ChildSize() > 0 {
				return append(res, t)
			}
		}
		if r.onlyOnce {
			break
		}
	}
	return nil
}

// Match 匹配
func (r RepeatParser) Match(lexer Lexer) bool {
	return r.parser.Match(lexer)
}

// Factory
type Factory struct {
	make0 func(arg interface{}) TreeNode // 通过闭包实现须重写的方法
}

// make
func (f *Factory) make(arg interface{}) TreeNode {
	return f.make0(arg)
}

// getFactory 获取Factory对象
func getFactory(treeType interface{}) *Factory {
	if treeType == nil {
		return nil
	}
	return &Factory{make0: func(arg interface{}) TreeNode {
		return NewTreeNode(arg)
	}}
}

// getForASTListFactory 获取Factory对象
func getForASTListFactory(treeType interface{}) *Factory {
	factory := getFactory(treeType)
	if factory == nil {
		factory = &Factory{make0: func(arg interface{}) TreeNode {
			switch arg.(type) {
			case []TreeNode:
				results := arg.([]TreeNode)
				if len(results) == 1 {
					return results[0]
				} else {
					return NewBranchNode(results)
				}
			}
			return nil
		}}
	}
	return factory
}

// Parser 解析器
type Parser struct {
	elements []ParserElement // 解析器元素
	factory  *Factory        //
}

// NewParser 创建Parser对象
func NewParser(treeType TreeNode) Parser {
	return Parser{
		elements: make([]ParserElement, 10),
		factory:  getForASTListFactory(treeType),
	}
}

// NewParserFromParser 创建Parser对象
func NewParserFromParser(parser Parser) Parser {
	return Parser{
		elements: parser.elements,
		factory:  parser.factory,
	}
}

// reset 重置解析器
func (p Parser) reset(treeType TreeNode) Parser {
	p.elements = make([]ParserElement, 10)
	p.factory = getForASTListFactory(treeType)
	return p
}

// Rule 获取解析器对象
func Rule() Parser {
	return RuleByType(nil)
}

// RuleByType 获取解析器对象
func RuleByType(treeType TreeNode) Parser {
	return NewParser(treeType)
}

// parse
func (p Parser) parse(lexer Lexer) TreeNode {
	result := make([]TreeNode, 0)
	for _, e := range p.elements {
		e.Parse(lexer, result)
	}
	return p.factory.make(result)
}

// Match
func (p Parser) Match(lexer Lexer) bool {
	if len(p.elements) == 0 {
		return true
	} else {
		e := p.elements[0]
		return e.Match(lexer)
	}
}

// Or or
func (p Parser) Or(parsers []Parser) Parser {
	p.elements = append(p.elements, NewOrTree(parsers))
	return p
}

// Sep
func (p Parser) Sep(pat ...string) Parser {
	p.elements = append(p.elements, NewSkip(pat))
	return p
}

// Ast
func (p Parser) Ast(parser Parser) Parser {
	p.elements = append(p.elements, NewTreeParser(parser))
	return p
}

// Number
func (p Parser) Number(typ interface{}) Parser {
	p.elements = append(p.elements, NewNumTokenParser(typ))
	return p
}

// Identifier
func (p Parser) Identifier(typ interface{}, r mapset.Set) Parser {
	p.elements = append(p.elements, NewIdTokenParser(typ, r))
	return p
}

// String
func (p Parser) String(typ interface{}) Parser {
	p.elements = append(p.elements, NewStrTokenParser(typ))
	return p
}

// Maybe
func (p Parser) Maybe(parser Parser) Parser {
	p2 := NewParserFromParser(parser)
	p2.reset(nil)
	p.elements = append(p.elements, NewOrTree([]Parser{parser, p2}))
	return p
}

// Option
func (p Parser) Option(parser Parser) Parser {
	p.elements = append(p.elements, NewRepeatParser(parser, true))
	return p
}

// Repeat
func (p Parser) Repeat(parser Parser) Parser {
	p.elements = append(p.elements, NewRepeatParser(parser, false))
	return p
}

// InsertChoice
func (p Parser) InsertChoice(parser Parser) Parser {
    e := p.elements[0]
	switch e.(type) {
	case OrTree:
		e.(OrTree).insert(parser)
		break
	default:
		otherwise := NewParserFromParser(p)
		p.reset(nil)
		p.Or([]Parser{parser, otherwise})
	}
    return p
}