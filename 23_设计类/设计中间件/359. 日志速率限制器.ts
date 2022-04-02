type TimeStamp = number
type Message = string

class Logger {
  private record: Map<Message, TimeStamp>

  constructor() {
    this.record = new Map()
  }

  // 每条 不重复 的消息最多只能每 10 秒打印一次
  shouldPrintMessage(timestamp: number, message: string): boolean {
    if (!this.record.has(message)) {
      this.record.set(message, timestamp)
      return true
    } else {
      if (this.record.get(message)! + 10 > timestamp) return false
      else {
        this.record.set(message, timestamp)
        return true
      }
    }
  }
}

// 类似于362
export {}
