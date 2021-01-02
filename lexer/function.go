package lexer

import "fmt"

// Function 函数定义对象
type Function struct {
	parameters ParameterListNode  // 参数列表
	body       BlockStatementNode // 函数体
	env        Environment        // 环境变量
}

// NewFunction 创建Function对象
func NewFunction(parameters ParameterListNode, body BlockStatementNode, env Environment) *Function {
	return &Function{
		parameters: parameters,
		body:       body,
		env:        env,
	}
}

// GetParameters 获取参数列表
func (f *Function) Parameters() ParameterListNode {
	return f.parameters
}

// Body 获取函数体
func (f *Function) Body() BlockStatementNode {
	return f.body
}

// makeEnv 获取环境变量
func (f *Function) makeEnv() Environment {
	return NewNestedEnvironment(f.env)
}

// String String方法
func (f *Function) String() string {
	return fmt.Sprintf("<fun: %v >", &f)
}
