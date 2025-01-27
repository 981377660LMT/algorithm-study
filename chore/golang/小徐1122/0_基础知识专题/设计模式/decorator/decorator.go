// • 继承强调的是等级制度和子类种类，这部分架构需要一开始就明确好
//
// • 装饰器模式强调的是“装饰”的过程，而不强调输入与输出，
//   能够动态地为对象增加某种特定的附属能力，相比于继承模式显得更加灵活，且符合开闭原则

package main

import (
	"context"
	"fmt"
)

type handleFunc func(ctx context.Context, param map[string]interface{}) error

func Decorate(fn handleFunc) handleFunc {
	return func(ctx context.Context, param map[string]interface{}) error {
		// 前处理
		fmt.Println("preprocess...")
		err := fn(ctx, param)
		fmt.Println("postprocess...")
		return err
	}
}
