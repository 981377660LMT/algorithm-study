type Id = number

interface Observable {
  subscribe: (subscriber: Subscriber) => Subscription
}

type Observer = Subscriber

interface Subscriber {
  next: (value: any) => void
  error: (error: Error) => void
  complete: () => void
}

interface Subscription {
  unsubscribe: () => void
}

interface Subject extends Observable, Observer {}

export type { Id, Observable, Observer, Subscriber, Subscription, Subject }
