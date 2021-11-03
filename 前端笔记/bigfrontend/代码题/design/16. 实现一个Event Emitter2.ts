type EventName = string
type Callback = (...args: any[]) => void

class EventEmitter {
  private events: Record<EventName, Callback[]>

  constructor() {
    this.events = {}
  }

  // 实现订阅
  on(type: EventName, callBack: Callback) {
    !this.events[type] && (this.events[type] = [])
    this.events[type].push(callBack)
  }

  // 删除订阅
  off(type: EventName, callBack: Callback) {
    if (!this.events[type]) return
    this.events[type] = this.events[type].filter(item => item !== callBack)
  }

  // 只执行一次订阅事件
  once(type: EventName, callBack: Callback) {
    const wrappedCallback = () => {
      callBack()
      this.off(type, wrappedCallback)
    }

    this.on(type, wrappedCallback)
  }

  // 触发事件
  emit(type: EventName, ...rest: any[]) {
    this.events[type] && this.events[type].forEach(cb => cb.call(this, ...rest))
  }
}

export {}
