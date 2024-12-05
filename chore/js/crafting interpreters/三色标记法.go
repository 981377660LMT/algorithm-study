package main

import "fmt"

func main() {
	// 创建对象图
	obj1 := &Object{}
	obj2 := &Object{}
	obj3 := &Object{}
	obj4 := &Object{}

	obj1.Children = []*Object{obj2, obj3}
	obj2.Children = []*Object{obj4}
	obj3.Children = []*Object{}
	obj4.Children = []*Object{}

	// 所有对象列表
	objects := []*Object{obj1, obj2, obj3, obj4}

	// 执行Mark阶段
	Mark(obj1)

	// 执行Sweep阶段
	objects = Sweep(objects)

	fmt.Printf("存活的对象数量: %d\n", len(objects))
}

type Color int

const (
	White Color = iota
	Gray
	Black
)

// Object 表示被管理的对象
type Object struct {
	Children []*Object
	Color    Color
}

// Mark 函数实现三色标记法
func Mark(obj *Object) {
	// 使用栈来避免递归
	stack := []*Object{obj}

	for len(stack) > 0 {
		// 弹出栈顶对象
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.Color == Black {
			continue
		}

		if current.Color == White {
			current.Color = Gray
			// 将当前对象重新压入栈中，以便之后标记为Black
			stack = append(stack, current)
			// 将所有子对象压入栈中
			for _, child := range current.Children {
				if child.Color == White {
					stack = append(stack, child)
				}
			}
		} else if current.Color == Gray {
			current.Color = Black
		}
	}
}

// Sweep 函数清理未标记的对象
func Sweep(objects []*Object) []*Object {
	var alive []*Object
	for _, obj := range objects {
		if obj.Color == Black {
			// 重置颜色以备下次GC
			obj.Color = White
			alive = append(alive, obj)
		} else {
			// 清理未被标记的对象，处理循环引用
			// 这里可以添加具体的清理逻辑，如释放资源
		}
	}
	return alive
}
