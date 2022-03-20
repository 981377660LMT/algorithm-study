interface Subscription {
  unsubscribe: () => void
}

interface IEventEmitter {
  on: (eventName: string, callback: (...args: any[]) => void) => Subscription
  emit: (eventName: string, ...args: any[]) => void
}

// 内部维护一个事件名到个个回调的map
/////////////////////////////////////////////////////////////////////////
interface Observable {
  subscribe: (subscriber: Subscriber) => Subscription
}

type Observer = Subscriber

interface Subscriber {
  next: (value: any) => void
  error: (error: Error) => void
  complete: () => void
}

// 内部维护一个id 到subscription 的 map
