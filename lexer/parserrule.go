package lexer

import (
	"fmt"
)

// ParserElement 解析器元素接口
type ParserElement interface {
	Parse(lexer Lexer, res []TreeNode) []TreeNode // 解析
	match(lexer Lexer) bool              // 匹配
}

// OrTree Or逻辑解析器元素
type OrTree struct {
	parsers []Parser
}

// NewOrTree 创建Or逻辑解析器元素
func NewOrTree(parsers []Parser) OrTree {
	return OrTree{parsers:parsers}
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

// match or逻辑匹配
func (o OrTree) match(lexer Lexer) bool {
	p := o.choose(lexer)
	return p != nil
}

// choose
func (o OrTree) choose(lexer Lexer) *Parser {
	for _, p := range o.parsers {
		if p.match(lexer) {
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
	factory *Factory //
}

// NewParser 创建Parser对象
func NewParser(treeType TreeNode) Parser {
	return Parser{
		elements: make([]ParserElement, 10),
		factory:  getForASTListFactory(treeType),
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

// match
func (p Parser) match(lexer Lexer) bool {
	if len(p.elements) == 0 {
		return true
	} else {
		e := p.elements[0]
		return e.match(lexer)
	}
}

// Or or
func (p Parser) Or(parsers []Parser) Parser {
	p.elements = append(p.elements, NewOrTree(parsers))
	return p
}
