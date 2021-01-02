package lexer

// Environment 环境对象接口
type Environment interface {
	Put(name string, value interface{})    // 保存对象
	PutNew(name string, value interface{}) // 添加新对象
	Get(name string) interface{}           // 获取值
	Where(name string) Environment         // 在所有作用域中获取值
}

// BasicEnvironment 基础环境对象实现
type BasicEnvironment struct {
	values map[string]interface{}
}

// NewBasicEnv 创建BasicEnvironment对象
func NewBasicEnv() BasicEnvironment {
	return BasicEnvironment{
		make(map[string]interface{}),
	}
}

// Put 保存对象
func (b BasicEnvironment) Put(name string, value interface{}) {
	b.values[name] = value
}

// PutNew 保存对象
func (b BasicEnvironment) PutNew(name string, value interface{}) {
	b.values[name] = value
}

// Get 获取值
func (b BasicEnvironment) Get(name string) interface{} {
	return b.values[name]
}

// Where 在所有作用域中获取值
func (b BasicEnvironment) Where(name string) Environment {
	return nil
}

const (
	TRUE  = 1
	FALSE = 0
)

// NestedEnvironment
type NestedEnvironment struct {
	values map[string]interface{} // 当前作用域变量
	outer  Environment            // 外层作用域变量
}

// NewNestedEnvironment 创建NestedEnvironment对象
func NewNestedEnvironment(environment Environment) NestedEnvironment {
	return NestedEnvironment{
		make(map[string]interface{}),
		environment,
	}
}

// SetOuter 设置外层变量
func (n NestedEnvironment) SetOuter(environment Environment) {
	n.outer = environment
}

// PutNew 保存新变量
func (n NestedEnvironment) PutNew(name string, value interface{}) {
	n.values[name] = value
}

// Where 在所有作用域中获取值
func (n NestedEnvironment) Where(name string) Environment {
	if n.values[name] != nil {
		return n
	}
	if n.outer == nil {
		return nil
	}
	return n.outer.Where(name)
}

// Put 保存对象
func (n NestedEnvironment) Put(name string, value interface{}) {
	e := n.Where(name)
	if e == nil {
		e = n
	}
	e.PutNew(name, value)
}

// Get 获取值
func (n NestedEnvironment) Get(name string) interface{} {
	v := n.values[name]
	if v == nil && n.outer != nil {
		return n.outer.Get(name)
	}
	return v
}
