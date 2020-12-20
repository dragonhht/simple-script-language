package lexer

import (
	"errors"
	"fmt"
	"strings"
)

// TreeNode 语法树节点
type TreeNode interface {
	Child(n int) (TreeNode, error) // 获取该节点下第n个子节点
	ChildSize() int                // 子节点个数
	Children() []TreeNode          // 获取子节点
	Location() string              // 定位显示
	String() string                // 实现String接口
}

// NewTreeNode 创建语法树节点
func NewTreeNode(treeType interface{}, arg interface{}) TreeNode {
	switch treeType.(type) {
	case PrimaryExpr:
		return CreatePrimaryExpr(arg.([]TreeNode))
	case NegativeExprNode:
		return NewNegativeExprNode(arg.([]TreeNode))
	case BlockStatementNode:
		return NewBlockStatementNode(arg.([]TreeNode))
	case NumberNode:
		return NewNumberNode(arg.(Token))
	case VariableNode:
		return NewVariableNode(arg.(Token))
	case StringNode:
		return NewStringNode(arg.(Token))
	case BinaryExprNode:
		return NewBinaryExprNode(arg.([]TreeNode))
	case IfStatementNode:
		return NewIfStatementNode(arg.([]TreeNode))
	case WhileStatementNode:
		return NewWhileStatementNode(arg.([]TreeNode))
	case NullStatementNode:
		return NewNullStatementNode(arg.([]TreeNode))
	}
	return nil
}

// LeafNode 语法树叶子节点
type LeafNode struct {
	token Token
	empty []TreeNode
}

// NewLeafNode 创建叶子节点
func NewLeafNode(token Token) LeafNode {
	return LeafNode{
		token: token,
		empty: make([]TreeNode, 0),
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
func (l LeafNode) Children() []TreeNode {
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

// Value 获取值
func (s StringNode) Value() string {
	return s.token.GetText()
}

// BranchNode 语法树树枝节点
type BranchNode struct {
	list []TreeNode
}

// NewBranchNode 创建树枝节点
func NewBranchNode(list []TreeNode) BranchNode {
	return BranchNode{list: list}
}

// Child 获取树枝节点下指定的子节点
func (b BranchNode) Child(n int) (TreeNode, error) {
	return b.list[n], nil
}

// ChildSize 子节点个数
func (b BranchNode) ChildSize() int {
	return len(b.list)
}

// Children 获取子节点
func (b BranchNode) Children() []TreeNode {
	return b.list
}

// Location 定位显示
func (b BranchNode) Location() string {
	for _, c := range b.list {
		s := c.Location()
		if s != "" {
			return s
		}
	}
	return ""
}

// String 实现String接口
func (b BranchNode) String() string {
	var buf strings.Builder
	buf.WriteString("(")
	sep := ""
	for _, v := range b.Children() {
		buf.WriteString(sep)
		sep = " "
		buf.WriteString(v.String())
	}
	buf.WriteString(")")
	return buf.String()
}

// NegativeExpr
type NegativeExprNode struct {
	BranchNode
}

// NewNegativeExprNode 创建NegativeExprNode对象
func NewNegativeExprNode(list []TreeNode) NegativeExprNode {
	return NegativeExprNode{
		NewBranchNode(list),
	}
}

// Operand
func (n NegativeExprNode) Operand() TreeNode {
	return n.list[0]
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
func NewBinaryExprNode(list []TreeNode) BinaryExprNode {
	return BinaryExprNode{
		NewBranchNode(list),
	}
}

// Left 获取子节点中的左子节点
func (b BinaryExprNode) Left() TreeNode {
	return b.list[0]
}

// Right 获取子节点中的右子节点
func (b BinaryExprNode) Right() TreeNode {
	return b.list[2]
}

// Operator 获取操作符
func (b BinaryExprNode) Operator() string {
	node := b.list[1]
	switch node.(type) {
	case LeafNode:
		return node.(LeafNode).token.GetText()
	}
	return ""
}

// PrimaryExpr
type PrimaryExpr struct {
	BranchNode
}

// PrimaryExpr
func NewPrimaryExpr(list []TreeNode) PrimaryExpr {
	return PrimaryExpr{NewBranchNode(list)}
}

// CreatePrimaryExpr
func CreatePrimaryExpr(list []TreeNode) TreeNode {
	if len(list) == 1 {
		return list[0]
	} else {
		return NewBranchNode(list)
	}
}

// BlockStatementNode
type BlockStatementNode struct {
	BranchNode
}

// NewBlockStatementNode
func NewBlockStatementNode(list []TreeNode) BlockStatementNode {
	return BlockStatementNode{NewBranchNode(list)}
}

// IfStatementNode
type IfStatementNode struct {
	BranchNode
}

// NewIfStatementNode
func NewIfStatementNode(list []TreeNode) IfStatementNode {
	return IfStatementNode{NewBranchNode(list)}
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
	return c
}

// ElseBlock else语句
func (i IfStatementNode) ElseBlock() TreeNode {
	if i.ChildSize() > 2 {
		c, err := i.Child(2)
		if err != nil {
			panic(err)
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
func NewWhileStatementNode(list []TreeNode) WhileStatementNode {
	return WhileStatementNode{NewBranchNode(list)}
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
func NewNullStatementNode(list []TreeNode) NullStatementNode {
	return NullStatementNode{NewBranchNode(list)}
}
