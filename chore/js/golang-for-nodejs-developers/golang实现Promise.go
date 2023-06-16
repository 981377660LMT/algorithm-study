// !golang实现Promise/golang手写Promise
// !Promise是一个容量为1的channel
//
// https://leexiaop.github.io/coding/4/ (实现Promise的链式调用)
// https://zhuanlan.zhihu.com/p/379740728
// https://juejin.cn/post/7083807190267461640
// https://juejin.cn/post/7110848876445171720
// https://zhuanlan.zhihu.com/p/61681036
// !resolve的入参可能是Promise 但是onFulfilled的入参不可能是Promise
//
// 大致思路:
// 1. Promise传一个executor函数, 该函数接收两个参数resolve和reject
// 2. Promise有三种状态: Pending, Fulfilled, Rejected
// 3. 调用wrappedResolve(value)后, Promise的状态变为Fulfilled, 并且将value存储在Promise的resolved字段中.
//    然后顺序取出successCallback中的函数并执行, 并将value作为参数传入.
//    调用wrappedReject(reason)后, Promise的状态变为Rejected, 并且将reason存储在Promise的rejected字段中.
//    然后顺序取出failureCallback中的函数并执行, 并将reason作为参数传入.
//    为了支持Await，内部还保存了一个容量为1的channel，当Promise的状态变为Fulfilled或Rejected时，会向该channel中写入数据。
// 4. Promise.Then() 链式调用: 传入两个函数onFulfilled和onRejected, 并返回一个新的Promise.
//    状态机,当检测到内部的Promise状态变为Fulfilled或Rejected时, 执行onFulfilled或onRejected.
//    如果onFulfilled或onRejected返回的是一个Promise, 则继续调用这个Promise的Then方法.
//    如果内部仍在Pending状态, 则将onFulfilled或onRejected存入successCallback或failureCallback中.
//    函数可以看成是惰性求值.
// 5. Promise.Await() 阻塞, 直到产生结果.
// 6. Promise.Catch() 等价于Promise.Then(nil, onRejected)
// 7. Promise.Finally() 等价于在中间加一层拦截.
// 8. Promise.Resolve 等价于NewPromise(func(resolve ResolveLike, reject RejectLike) { resolve(value) })
// 9. Promise.Reject 等价于NewPromise(func(resolve ResolveLike, reject RejectLike) { reject(reason) })
// 10. 还有一系列的静态方法, 如Promise.All, Promise.Race等等, 这些可以通过遍历任务列表来封装实现.

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	NewPromise(func(resolve ResolveLike, reject RejectLike) {
		time.Sleep(500 * time.Millisecond)
		if rand := rand.Intn(10); rand > 5 {
			resolve("ok")
		} else {
			reject(errors.New("failed"))
		}
	}).Then(func(value interface{}) interface{} {
		fmt.Println(value)
		return 0
	}, func(reason error) {
		fmt.Println(reason)
	})

	p2 := NewPromise(func(resolve ResolveLike, _ RejectLike) {
		time.Sleep(100 * time.Millisecond)
		resolve(100)
	}).Await()
	fmt.Println(p2)

	time.Sleep(1 * time.Second)

}

func NewPromise(
	executor func(resolve ResolveLike, reject RejectLike),
) *Promise {
	res := &Promise{state: Pending, queue: make(chan interface{}, 1)}
	wrappedResolve := func(value interface{}) interface{} {
		if res.state != Pending {
			return nil
		}
		res.state = Fulfilled
		res.resolved = value
		for len(res.successCallback) > 0 {
			res.successCallback[0](value)
			res.successCallback = res.successCallback[1:]
		}
		res.queue <- value
		return nil
	}

	wrappedReject := func(reason error) {
		if res.state != Pending {
			return
		}
		res.state = Rejected
		res.rejected = reason
		for len(res.failureCallback) > 0 {
			res.failureCallback[0](reason)
			res.failureCallback = res.failureCallback[1:]
		}
		res.queue <- reason
	}

	executor(wrappedResolve, wrappedReject)

	return res
}

type PromiseState uint8

// resolve的入参可能是Promise,但是onFulfilled的入参不能是Promise.
type ResolveLike = func(value interface{}) interface{}
type RejectLike = func(reason error)

const (
	Pending   PromiseState = 0
	Fulfilled PromiseState = 1
	Rejected  PromiseState = 2
)

type Promise struct {
	state           PromiseState
	resolved        interface{}
	rejected        error
	successCallback []ResolveLike
	failureCallback []RejectLike
	queue           chan interface{} // support await
}

func (p *Promise) Await() interface{} {
	return <-p.queue
}

func (p *Promise) Then(
	onFulfilled ResolveLike,
	onRejected RejectLike,
) *Promise {
	// 使得onFulfilled可以处理Promise
	wrap := func(onFulfilled ResolveLike, resolve ResolveLike, reject RejectLike) ResolveLike {
		return func(res interface{}) interface{} {
			if promise, ok := res.(*Promise); ok {
				promise.Then(resolve, reject)
			} else {
				onFulfilled(res)
			}
			return nil
		}
	}

	return NewPromise(func(resolve ResolveLike, reject RejectLike) {
		if p.state == Fulfilled {
			res := onFulfilled(p.resolved)
			if promise, ok := res.(*Promise); ok {
				promise.Then(resolve, reject)
			} else {
				resolve(res)
			}
		} else if p.state == Rejected {
			onRejected(p.rejected)
		} else {
			p.successCallback = append(p.successCallback, wrap(onFulfilled, resolve, reject))
			p.failureCallback = append(p.failureCallback, onRejected)
		}
	})
}

func (p *Promise) Finally(onFinally func()) *Promise {
	return NewPromise(func(resolve ResolveLike, reject RejectLike) {
		p.Then(
			func(value interface{}) interface{} {
				onFinally()
				resolve(value)
				return nil
			},
			func(reason error) {
				onFinally()
				reject(reason)
			})
	})
}

func NewPromiseResolve(value interface{}) *Promise {
	if promise, ok := value.(*Promise); ok {
		return promise
	}
	return NewPromise(func(resolve ResolveLike, _ RejectLike) {
		resolve(value)
	})
}

func NewPromiseReject(reason error) *Promise {
	return NewPromise(func(_ ResolveLike, reject RejectLike) {
		reject(reason)
	})
}
