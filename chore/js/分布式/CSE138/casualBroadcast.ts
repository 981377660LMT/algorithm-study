interface Message {
  content: string
  senderId: number
  vectorClock: number[]
}

class VectorClock {
  private readonly _clock: number[]

  constructor(size: number) {
    this._clock = Array(size).fill(0)
  }

  increment(index: number): void {
    if (index >= 0 && index < this._clock.length) {
      this._clock[index]++
    }
  }

  update(msgClock: number[]): void {
    for (let i = 0; i < this._clock.length; i++) {
      this._clock[i] = Math.max(this._clock[i], msgClock[i])
    }
  }

  /**
   * 检查是否可以交付消息
   * 确保消息按因果顺序交付。
   */
  canDeliver(msgClock: number[], senderId: number): boolean {
    for (let i = 0; i < this._clock.length; i++) {
      if (i === senderId) {
        if (this._clock[i] + 1 !== msgClock[i]) return false
      } else {
        if (this._clock[i] < msgClock[i]) return false
      }
    }
    return true
  }

  // 复制当前时钟
  copy(): number[] {
    return [...this._clock]
  }

  toString(): string {
    return `[${this._clock.join(',')}]`
  }
}

class Node {
  private readonly _nodeId: number
  private readonly _vectorClock: VectorClock
  private readonly _messageBuffer: Message[]
  private readonly _delivered: Set<string>
  private readonly _nodes: Node[]

  constructor(nodeId: number, numNodes: number) {
    this._nodeId = nodeId
    this._vectorClock = new VectorClock(numNodes)
    this._messageBuffer = []
    this._delivered = new Set()
    this._nodes = []
  }

  addNode(node: Node): void {
    this._nodes.push(node)
  }

  broadcast(content: string): void {
    // 增加自己的向量时钟
    this._vectorClock.increment(this._nodeId)

    const message: Message = {
      content,
      senderId: this._nodeId,
      vectorClock: this._vectorClock.copy()
    }

    console.log(
      `Node ${this._nodeId} broadcasting: "${content}" with clock ${this._vectorClock.toString()}`
    )

    // 发送给所有其他节点
    this._nodes.forEach(node => node.receive(message))
  }

  // 接收消息
  receive(message: Message): void {
    // 生成消息唯一标识
    const msgId = `${message.senderId}-${message.content}-${message.vectorClock.join(',')}`

    // 如果消息已经交付过，直接返回
    if (this._delivered.has(msgId)) {
      return
    }

    console.log(
      `Node ${this._nodeId} received: "${message.content}" from ${message.senderId} with clock ${message.vectorClock}`
    )

    if (this._vectorClock.canDeliver(message.vectorClock, message.senderId)) {
      this._deliver(message)
      this._tryDeliverBuffered()
    } else {
      console.log(`Node ${this._nodeId} buffering message: "${message.content}"`)
      this._messageBuffer.push(message)
    }
  }

  // 获取当前向量时钟
  getVectorClock(): string {
    return this._vectorClock.toString()
  }

  // 交付消息
  private _deliver(message: Message): void {
    const msgId = `${message.senderId}-${message.content}-${message.vectorClock.join(',')}`

    console.log(`Node ${this._nodeId} delivering: "${message.content}"`)

    this._vectorClock.update(message.vectorClock)
    this._delivered.add(msgId)
  }

  private _tryDeliverBuffered(): void {
    let delivered = true
    while (delivered) {
      delivered = false
      const toDeliver: Message[] = []

      // 收集所有可以交付的消息
      for (let i = 0; i < this._messageBuffer.length; i++) {
        if (
          this._vectorClock.canDeliver(
            this._messageBuffer[i].vectorClock,
            this._messageBuffer[i].senderId
          )
        ) {
          toDeliver.push(this._messageBuffer[i])
        }
      }

      // 交付并移除可以交付的消息
      if (toDeliver.length > 0) {
        delivered = true
        for (const message of toDeliver) {
          const index = this._messageBuffer.indexOf(message)
          if (index !== -1) {
            this._messageBuffer.splice(index, 1)
            this._deliver(message)
          }
        }
      }
    }
  }
}

// 测试代码
function main() {
  // 创建三个节点
  const numNodes = 3
  const nodes = Array.from({ length: numNodes }, (_, i) => new Node(i, numNodes))

  // 设置节点之间的连接
  nodes.forEach(node => {
    nodes.forEach(other => {
      if (node !== other) {
        node.addNode(other)
      }
    })
  })

  console.log('=== Starting Causal Broadcast Test ===')

  // 节点0广播消息A
  nodes[0].broadcast('A')

  // 节点1在收到A后广播消息B
  nodes[1].broadcast('B')

  // 节点2在收到A和B后广播消息C
  nodes[2].broadcast('C')

  console.log('\n=== Final Vector Clocks ===')
  nodes.forEach((node, i) => {
    console.log(`Node ${i}: ${node.getVectorClock()}`)
  })
}

// 运行测试
main()

export {}
