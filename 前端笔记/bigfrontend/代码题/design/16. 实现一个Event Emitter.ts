import { Func } from '../typings'

// class Subscription {
//   constructor(
//     private subscriber: Map<EventName, Func[]>,
//     private eventName: EventName,
//     private callback: Func
//   ) {}

//   release() {
//     if (this.subscriber.has(this.eventName)) {

//     }
//   }
// }

type EventName = string
type Id = number

class EventEmitter {
  private cbRecord: Record<EventName, Record<Id, Func>>
  private uuid: number

  constructor() {
    this.cbRecord = {}
    this.uuid = 0
  }

  subscribe(eventName: string, callback: Func) {
    const uuid = this.uuid++
    if (!(eventName in this.cbRecord)) this.cbRecord[eventName] = {}
    this.cbRecord[eventName][uuid] = callback

    return {
      release: () => {
        delete this.cbRecord[eventName][uuid]
      },
    }
  }

  emit(eventName: string, ...args: any[]): void {
    if (this.cbRecord[eventName]) {
      Object.values(this.cbRecord[eventName]).forEach(cb => cb.call(null, ...args))
    }
  }

  static main() {
    const cb1 = () => console.log(1)
    const cb2 = () => console.log(2)
    const emitter = new EventEmitter()
    const sub1 = emitter.subscribe('event1', cb1)
    const sub2 = emitter.subscribe('event2', cb2)
    // 同一个callback可以重复订阅同一个事件
    const sub3 = emitter.subscribe('event1', cb1)
    // callback1 会被调用两次
    emitter.emit('event1', 1, 2)
    sub1.release()
    sub3.release()
    // 现在即使'event1'被触发,
    // callback1 也不会被调用
    emitter.emit('event1', 1, 2)
  }
}

EventEmitter.main()
export {}
