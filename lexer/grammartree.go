package lexer

import (
	"errors"
	"fmt"
	"simple-script-language/utils/list"
	"strings"
)

// TreeNode 语法树节点
type TreeNode interface {
	Child(n int) (TreeNode, error) // 获取该节点下第n个子节点
	ChildSize() int                // 子节点个数
	Children() *list.ArrayList     // 获取子节点
	Location() string              // 定位显示
	String() string                // 实现String接口
	Eval(environment Environment) interface{} // 获取节点计算值
}

// NewTreeNode 创建语法树节点
func NewTreeNode(treeType interface{}, arg interface{}) TreeNode {
	switch treeType.(type) {
	case PrimaryExpr:
		return CreatePrimaryExpr(arg.(*list.ArrayList))
	case NegativeExprNode:
		return NewNegativeExprNode(arg.(*list.ArrayList))
	case BlockStatementNode:
		return NewBlockStatementNode(arg.(*list.ArrayList))
	case NumberNode:
		return NewNumberNode(arg.(Token))
	case VariableNode:
		return NewVariableNode(arg.(Token))
	case StringNode:
		return NewStringNode(arg.(Token))
	case BinaryExprNode:
		return NewBinaryExprNode(arg.(*list.ArrayList))
	case IfStatementNode:
		return NewIfStatementNode(arg.(*list.ArrayList))
	case WhileStatementNode:
		return NewWhileStatementNode(arg.(*list.ArrayList))
	case NullStatementNode:
		return NewNullStatementNode(arg.(*list.ArrayList))
	}
	return nil
}

// LeafNode 语法树叶子节点
type LeafNode struct {
	token Token
	empty *list.ArrayList
}

// NewLeafNode 创建叶子节点
func NewLeafNode(token Token) LeafNode {
	return LeafNode{
		token: token,
		empty: list.New(0),
	}
}

// Child 获取叶子节点下指定的子节点(因叶子节点没有子节点，则调用会报错)
func (l LeafNode) Child(n int) (TreeNode, error) {
	return nil, errors.New("叶子节点不存在子节点")
}

// ChildSize 子节点个数
func (l LeafNode) ChildSize() int {
	return 0
}

// Children 获取子节点
func (l LeafNode) Children() *list.ArrayList {
	return l.empty
}

// Location 定位显示
func (l LeafNode) Location() string {
	return fmt.Sprintf("at line %v", l.token.GetLineNumber())
}

// String 实现String方法方便打印
func (l LeafNode) String() string {
	return l.token.GetText()
}

func (l LeafNode) Eval(env Environment) interface{} {
	panic(fmt.Sprintf("cannot eval: %v", l.String()))
}

// NumberNode 数值型叶子节点
type NumberNode struct {
	LeafNode
}

// NewNumberNode 创建NumberNode对象
func NewNumberNode(token Token) NumberNode {
	return NumberNode{
		LeafNode: NewLeafNode(token),
	}
}

// Eval 获取计算值
func (n NumberNode) Eval(env Environment) interface{} {
	return n.Value()
}

// Value 获取值
func (n NumberNode) Value() int {
	num, _ := n.token.GetNumber()
	return num
}

// VariableNode 变量叶子节点
type VariableNode struct {
	LeafNode
}

// NewVariableNode 创建VariableNode对象
func NewVariableNode(token Token) VariableNode {
	return VariableNode{LeafNode: NewLeafNode(token)}
}

// Eval 获取计算值
func (v VariableNode) Eval(env Environment) interface{} {
	value := env.Get(v.Name())
	if value == nil {
		panic(fmt.Sprintf("undefined name: %v", v.Name()))
	}
	return value
}

// Name 获取变量名
func (v VariableNode) Name() string {
	return v.token.GetText()
}

// StringNode
type StringNode struct {
	LeafNode
}

// NewStringNode 创建StringNode对象
func NewStringNode(token Token) StringNode {
	return StringNode{LeafNode: NewLeafNode(token)}
}

// Eval 获取计算值
func (s StringNode) Eval(env Environment) interface{} {
	return s.Value()
}

// Value 获取值
func (s StringNode) Value() string {
	return s.token.GetText()
}

// BranchNode 语法树树枝节点
type BranchNode struct {
	list *list.ArrayList
}

// NewBranchNode 创建树枝节点
func NewBranchNode(list *list.ArrayList) BranchNode {
	return BranchNode{list: list}
}

// Eval 获取计算值
func (b BranchNode) Eval(env Environment) interface{} {
	panic(fmt.Sprintf("cannot eval: %v", b.String()))
}

// Child 获取树枝节点下指定的子节点
func (b BranchNode) Child(n int) (TreeNode, error) {
	node, err := b.list.Get(n)
	return node.(TreeNode), err
}

// ChildSize 子节点个数
func (b BranchNode) ChildSize() int {
	return b.list.Size()
}

// Children 获取子节点
func (b BranchNode) Children() *list.ArrayList {
	return b.list
}

// Location 定位显示
func (b BranchNode) Location() string {
	result := ""
	b.list.For(func(k int, v interface{}) {
		c := v.(TreeNode)
		s := c.Location()
		if s != "" {
			result = s
			return
		}
	})
	return result
}

// String 实现String接口
func (b BranchNode) String() string {
	var buf strings.Builder
	buf.WriteString("(")
	sep := ""
	b.Children().For(func(k int, v interface{}) {
		buf.WriteString(sep)
		sep = " "
		buf.WriteString(v.(TreeNode).String())
	})
	buf.WriteString(")")
	return buf.String()
}

// NegativeExpr
type NegativeExprNode struct {
	BranchNode
}

// NewNegativeExprNode 创建NegativeExprNode对象
func NewNegativeExprNode(list *list.ArrayList) NegativeExprNode {
	return NegativeExprNode{
		NewBranchNode(list),
	}
}

// Eval 获取计算值
func (n NegativeExprNode) Eval(env Environment) interface{} {
	value := n.Operand().Eval(env)
	switch value.(type) {
	case int:
		return -value.(int)
	}
	panic(fmt.Sprintf("bad type for -"))
}

// Operand
func (n NegativeExprNode) Operand() TreeNode {
	node, _ := n.list.Get(0)
	return node.(TreeNode)
}

// String
func (n NegativeExprNode) String() string {
	return fmt.Sprintf("-%v", n.Operand())
}

// BinaryExprNode 双目运算表达式节点
type BinaryExprNode struct {
	BranchNode
}

// NewBinaryExprNode 创建BinaryExprNode对象
func NewBinaryExprNode(list *list.ArrayList) BinaryExprNode {
	return BinaryExprNode{
		NewBranchNode(list),
	}
}

// Eval 获取计算值
func (b BinaryExprNode) Eval(env Environment) interface{} {
	op := b.Operator()
	if op == "=" {
		right := b.Right().Eval(env)
		return b.computeAssign(env, right)
	}
	left := b.Left().Eval(env)
	right := b.Right().Eval(env)
	return b.computeOp(left, op, right)
}

// Left 获取子节点中的左子节点
func (b BinaryExprNode) Left() TreeNode {
	node, _ := b.list.Get(0)
	return node.(TreeNode)
}

// Right 获取子节点中的右子节点
func (b BinaryExprNode) Right() TreeNode {
	node, _ := b.list.Get(2)
	return node.(TreeNode)
}

// Operator 获取操作符
func (b BinaryExprNode) Operator() string {
	node, _ := b.list.Get(1)
	switch node.(type) {
	case LeafNode:
		return node.(LeafNode).token.GetText()
	}
	return ""
}

// computeAssign 表达式复制操作
func (b BinaryExprNode) computeAssign(env Environment, rightVal interface{}) interface{} {
	left := b.Left()
	switch left.(type) {
	case VariableNode:
		env.Put(left.(VariableNode).Name(), rightVal)
		return rightVal
	}
	panic(fmt.Sprintf("bad assignment"))
}

// computeOp 表达式计算
func (b BinaryExprNode) computeOp(left interface{}, op string, right interface{}) interface{} {
	nl, lok := left.(int)
	nr, rok := right.(int)
	if lok && rok {
		return computeNumber(nl, op, nr)
	}
	if op == "+" {
		return fmt.Sprintf("%s%s", left, right)
	}
	if op == "==" {
		if left == nil {
			if right == nil {
				return TRUE
			} else {
				return FALSE
			}
		}
		if left == right {
			return TRUE
		} else {
			return FALSE
		}
	}
	panic(fmt.Sprintf("bad type"))
}

// computeNumber 整型计算
func computeNumber(left int, op string, right int) interface{} {
	switch op {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	case "%":
		return left % right
	case "==":
		if left == right {
			return TRUE
		} else {
			return FALSE
		}
	case ">":
		if left > right {
			return TRUE
		} else {
			return FALSE
		}
	case "<":
		if left < right {
			return TRUE
		} else {
			return FALSE
		}
	}
	panic(fmt.Sprintf("bad operator"))
}

// PrimaryExpr
type PrimaryExpr struct {
	BranchNode
}

// PrimaryExpr
func NewPrimaryExpr(list *list.ArrayList) PrimaryExpr {
	return PrimaryExpr{NewBranchNode(list)}
}

// CreatePrimaryExpr
func CreatePrimaryExpr(list *list.ArrayList) TreeNode {
	if list.Size() == 1 {
		node, _ := list.Get(0)
		return node.(TreeNode)
	} else {
		return NewBranchNode(list)
	}
}

// BlockStatementNode
type BlockStatementNode struct {
	BranchNode
}

// NewBlockStatementNode
func NewBlockStatementNode(list *list.ArrayList) BlockStatementNode {
	return BlockStatementNode{NewBranchNode(list)}
}

// Eval 获取计算值
func (b BlockStatementNode) Eval(env Environment) interface{} {
	var result interface{}
	result = 0
	b.Children().For(func(k int, v interface{}) {
		_, ok := v.(NullStatementNode)
		if !ok {
			result = v.(TreeNode).Eval(env)
		}
	})
	return result
}

// IfStatementNode
type IfStatementNode struct {
	BranchNode
}

// NewIfStatementNode
func NewIfStatementNode(list *list.ArrayList) IfStatementNode {
	if list.Size() > 0 {
		first, _ := list.Get(0)
		if fv,fok := first.(LeafNode); fok && fv.token.GetText() == "if" {
			list.Remove(0)
		}
	}
	return IfStatementNode{NewBranchNode(list)}
}

// Eval 获取计算值
func (i IfStatementNode) Eval(env Environment) interface{} {
	c := i.Condition().Eval(env)
	cv, cok := c.(int)
	if cok && cv != FALSE {
		return i.ThenBlock().Eval(env)
	}
	b := i.ElseBlock()
	if b == nil {
		return 0
	}
	return b.Eval(env)
}

// Condition 条件
func (i IfStatementNode) Condition() TreeNode {
	c, err := i.Child(0)
	if err != nil {
		panic(err)
	}
	return c
}

// ThenBlock 条件为真时的语句
func (i IfStatementNode) ThenBlock() TreeNode {
	c, err := i.Child(1)
	if err != nil {
		panic(err)
	}
	if c.ChildSize() > 0 {
		first, _ := c.Child(0)
		if cv, cok := first.(LeafNode); cok && cv.token.GetText() == "{" {
			c.Children().Remove(0)
		}
	}
	if c.ChildSize() > 0 {
		last, _ := c.Child(c.ChildSize() - 1)
		if cv, cok := last.(LeafNode); cok && cv.token.GetText() == "}" {
			c.Children().Remove(c.ChildSize() - 1)
		}
	}
	return c
}

// ElseBlock else语句
func (i IfStatementNode) ElseBlock() TreeNode {
	if i.ChildSize() > 2 {
		c, err := i.Child(2)
		if err != nil {
			panic(err)
		}
		first, ferr := c.Child(0)
		if ferr == nil {
			if fv,fok := first.(LeafNode); fok && fv.token.GetText() == "else"  {
				c.Children().Remove(0)
			}
		}
		if c.ChildSize() > 0 {
			block, _ := c.Child(0)
			if block.ChildSize() > 0 {
				first, _ := block.Child(0)
				if cv, cok := first.(LeafNode); cok && cv.token.GetText() == "{" {
					block.Children().Remove(0)
				}
			}
			if block.ChildSize() > 0 {
				last, _ := block.Child(block.ChildSize() - 1)
				if cv, cok := last.(LeafNode); cok && cv.token.GetText() == "}" {
					block.Children().Remove(block.ChildSize() - 1)
				}
			}
		}
		return c
	}
	return nil
}

// String
func (i IfStatementNode) String() string {
	return fmt.Sprintf("(if %v %v else %v)", i.Condition(), i.ThenBlock(), i.ElseBlock())
}

// WhileStatementNode
type WhileStatementNode struct {
	BranchNode
}

// NewWhileStatementNode
func NewWhileStatementNode(list *list.ArrayList) WhileStatementNode {
	if list.Size() > 0 {
		first, _ := list.Get(0)
		if fv,fok := first.(LeafNode); fok && fv.token.GetText() == "while" {
			list.Remove(0)
		}
	}
	return WhileStatementNode{NewBranchNode(list)}
}

// Eval 获取计算值
func (w WhileStatementNode) Eval(env Environment) interface{} {
	var result interface{}
	result = 0
	for {
		c := w.Condition().Eval(env)
		cv, cok := c.(int)
		if cok && cv == FALSE {
			return result
		}
		result = w.Body().Eval(env)
	}
}

// Condition 条件
func (w WhileStatementNode) Condition() TreeNode {
	c, err := w.Child(0)
	if err != nil {
		panic(err)
	}
	return c
}

// Body 条件为真时的语句
func (w WhileStatementNode) Body() TreeNode {
	c, err := w.Child(1)
	if err != nil {
		panic(err)
	}
	if c.ChildSize() > 0 {
		first, _ := c.Child(0)
		if cv, cok := first.(LeafNode); cok && cv.token.GetText() == "{" {
			c.Children().Remove(0)
		}
	}
	if c.ChildSize() > 0 {
		last, _ := c.Child(c.ChildSize() - 1)
		if cv, cok := last.(LeafNode); cok && cv.token.GetText() == "}" {
			c.Children().Remove(c.ChildSize() - 1)
		}
	}
	return c
}

// String
func (w WhileStatementNode) String() string {
	return fmt.Sprintf("(while %v %v)", w.Condition(), w.Body())
}

// NullStatementNode
type NullStatementNode struct {
	BranchNode
}

// NewNullStatementNode
func NewNullStatementNode(list *list.ArrayList) NullStatementNode {
	return NullStatementNode{NewBranchNode(list)}
}
