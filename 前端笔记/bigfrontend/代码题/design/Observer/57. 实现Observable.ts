// 你用过RxJS吗? 其中最重要的概念就是Observable 和 Observer。
// Observable 决定了values如何传递给Observer，Observer本质上就是一些callback的集合。

import type { Func } from '../../typings'
import type { Id, Subscriber, Subscription } from './typings'

class Observable {
  private setup: (subscriber: Subscriber) => void
  private id: Id
  private subRecord: Map<Id, Subscription> // 支持多个订阅

  constructor(setup: (subscriber: Subscriber) => void) {
    this.setup = setup
    this.id = 0
    this.subRecord = new Map()
  }

  subscribe(subscriber: Subscriber | Func): Subscription {
    const observable = this
    const subscriberWrapper = {
      id: this.id++,
      isSubscribed: true,
      next(value: any) {
        if (!this.isSubscribed) return
        if (subscriber instanceof Function) return subscriber(value)
        else return subscriber.next ? subscriber.next(value) : null
      },
      error(error: Error) {
        if (!this.isSubscribed) return
        this.unsubscribe()
        if (subscriber instanceof Function) return
        return subscriber.error ? subscriber.error(error) : null
      },
      complete() {
        if (!this.isSubscribed) return
        this.unsubscribe()
        if (subscriber instanceof Function) return
        return subscriber.complete ? subscriber.complete() : null
      },
      unsubscribe() {
        if (!this.isSubscribed) return
        this.isSubscribed = false
        observable.subRecord.delete(this.id)
      },
    }

    this.subRecord.set(subscriberWrapper.id, subscriberWrapper)
    this.setup(subscriberWrapper)
    return subscriberWrapper
  }
}

if (require.main === module) {
  // 观察者即订阅者,很明显就是3个callback而已。
  const observer = {
    next: (value: any) => {
      console.log('we got a value', value)
    },
    error: (error: Error) => {
      console.log('we got an error', error)
    },
    complete: () => {
      console.log('ok, no more values')
    },
  }

  // 可以把让这个Observer订阅一个Observable，Observable 会传递给这个Observer以值(value)或者error。
  // setup为接下来observable.subscribe的每一个订阅者进行初始化操作
  const observable = new Observable(subscriber => {
    subscriber.next(1)
    subscriber.next(2)
    setTimeout(() => {
      subscriber.next(3)
      subscriber.next(4)
      subscriber.complete()
    }, 100)
  })

  // subscribe() 返回的是一个Subscription ，这个subscription可以用来取消订阅。
  const subscription = observable.subscribe(observer)
  setTimeout(() => {
    subscription.unsubscribe()
  }, 50)
  // 1.error 和 complete只能触发一次。其后的next/error/complete 需要被忽略。(即触发后立即取消订阅)
  // 2.在订阅的时候next/error/complete需要都不是必须。如果传入的是一个函数，这个函数需要被默认为next。
  // 3.需要支持多个订阅。

  // 总结:
  // 1.两个er是什么? observer和subscriber是同一个东西(视角不同名字就不一样)，即带有几个回调的对象
  // 2.observerble是什么？订阅中心;observalble初始化setup，observalble.subscribe时触发setup并返回一个subscription
  // 3.subscription是什么？一个带有unsubscribe的对象
  // 4.observerble如何保存订阅信息？使用一个map记录订阅者id到subscription的对应关系
}

export { Observable }
