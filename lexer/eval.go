package lexer

// Environment 环境对象接口
type Environment interface {
	Put(name string, value interface{}) // 保存对象
	Get(name string) interface{}        // 获取值
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

// Get 获取值
func (b BasicEnvironment) Get(name string) interface{} {
	return b.values[name]
}

const (
	TRUE  = 1
	FALSE = 0
)
