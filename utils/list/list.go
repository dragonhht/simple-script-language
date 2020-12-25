package list

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// ArrayList 数组列表
type ArrayList struct {
	list []interface{} // 内部数组
	size int           // 长度
}

// New 创建数组列表
func New(initLen int) *ArrayList {
	return &ArrayList{
		list: make([]interface{}, initLen),
		size: 0,
	}
}

// Get 获取指定索引的值
func (a *ArrayList) Get(index int) (interface{}, error) {
	if index < 0 || index > a.size {
		return nil, errors.New(fmt.Sprintf("ArrayIndexOutOfBounds: size: %v, index: %v", a.size, index))
	}
	return a.list[index], nil
}

// add 添加
func (a *ArrayList) Add(item interface{}) {
	if a.size < len(a.list)-1 {
		a.list[a.size] = item
		a.size++
		return
	}
	a.list = expansion(a.list)
	a.list[a.size] = item
	a.size++
}

// expansion 数组扩容
func expansion(list []interface{}) []interface{} {
	len := len(list)
	newList := make([]interface{}, len+(len)/2)
	copy(newList, list)
	return newList
}

// Remove 移除指定的数据
func (a *ArrayList) Remove(index int) (interface{}, error) {
	if index < 0 || index > a.size {
		return nil, errors.New(fmt.Sprintf("ArrayIndexOutOfBounds: size: %v, index: %v", a.size, index))
	}
	listLen := len(a.list)
	item := a.list[index]
	var newList []interface{}
	if a.size-1 < listLen/2 {
		newList = make([]interface{}, listLen/2)
	} else {
		newList = make([]interface{}, listLen)
	}
	list1 := a.list[0:index]
	list2 := a.list[index+1 : a.size]
	copy(newList, list1)
	copy(newList[len(list1):], list2)
	a.list = newList
	a.size--
	return item, nil
}

// Clear 清空
func (a *ArrayList) Clear() {
	a.list = make([]interface{}, 10)
	a.size = 0
}

// String
func (a *ArrayList) String() string {
	var buf bytes.Buffer
	index := 0
	for _, v := range a.list {
		if index < a.size {
			buf.WriteString(fmt.Sprintf("%v, ", v))
			index++
		} else {
			break
		}
	}
	str := buf.String()
	if strings.HasSuffix(str, ", ") {
		str = str[0 : len(str)-2]
	}
	return str
}

// For 循环
func (a *ArrayList) For(handler func(k int, v interface{})) {
	for lk, lv := range a.list {
		if lk < a.size {
			handler(lk, lv)
		} else {
			break
		}
	}
}

// Size 大小
func (a *ArrayList) Size() int {
	return a.size
}
