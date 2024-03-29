type Id = number

/**
 * @description 可观察对象，注册中心
 */
interface Observable {
  subscribe: (subscriber: Subscriber) => Subscription // !addEventListener
}

type Observer = Subscriber

interface Subscriber {
  // !callback
  next: (value: any) => void
  error: (error: Error) => void
  complete: () => void
}

interface Subscription {
  unsubscribe: () => void // !removeEventListener
}

/**
 * @description 注册中心(Observable)，并且可以广播(Observer)
 */
interface Subject extends Observable, Observer {}

export type { Id, Observable, Observer, Subscriber, Subscription, Subject }
