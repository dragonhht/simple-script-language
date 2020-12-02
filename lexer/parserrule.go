package lexer

// ParserElement 解析器元素接口
type ParserElement interface {
	Parse(lexer Lexer, res []TreeNode) // 解析
	match(lexer Lexer)                 // 匹配
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

// Parser 解析器
type Parser struct {
	elements []ParserElement // 解析器元素

}
