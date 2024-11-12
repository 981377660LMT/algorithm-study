/* eslint-disable no-lonely-if */
// ChandyLamport.ts

interface Message {
  type: 'normal' | 'marker'
  content: string
}

class Channel {
  readonly fromProcess: Process
  readonly toProcess: Process
  private readonly _queue: Message[] = []

  constructor(fromProcess: Process, toProcess: Process) {
    this.fromProcess = fromProcess
    this.toProcess = toProcess
  }

  send(message: Message) {
    this._queue.push(message)
  }

  receive(): Message | undefined {
    return this._queue.shift()
  }
}

class Process {
  private readonly _processId: number
  private readonly _incomingChannels: Channel[] = []
  private readonly _outgoingChannels: Channel[] = []
  private _isRecording = false
  private _recordedState: any = undefined
  private _recordedMessages: Map<Channel, Message[]> = new Map()
  private _receivedMarkers: Set<Channel> = new Set()

  constructor(processId: number) {
    this._processId = processId
  }

  addIncomingChannel(channel: Channel) {
    this._incomingChannels.push(channel)
  }

  addOutgoingChannel(channel: Channel) {
    this._outgoingChannels.push(channel)
  }

  // 模拟进程运行，检查消息队列
  run() {
    setInterval(() => {
      for (const channel of this._incomingChannels) {
        const message = channel.receive()
        if (message) {
          this._receiveMessage(channel, message)
        }
      }
    }, 100)
  }

  sendMessage(targetProcess: Process, content: string) {
    const channel = this._outgoingChannels.find(ch => ch.toProcess === targetProcess)
    if (channel) {
      const message: Message = { type: 'normal', content }
      channel.send(message)
      console.log(
        `Process ${this._processId} sent message to Process ${targetProcess._processId}:`,
        content
      )
    }
  }

  // 发起快照
  startSnapshot() {
    if (!this._isRecording) {
      this._isRecording = true
      this._recordedState = this._getState()
      console.log(
        `Process ${this._processId} initiates snapshot and records its state:`,
        this._recordedState
      )
      this._sendMarker()
      // 开始记录所有输入通道的消息
      for (const ch of this._incomingChannels) {
        this._recordedMessages.set(ch, [])
      }
    }
  }

  // 获取记录的快照
  getSnapshot() {
    const snapshot: any = {
      processId: this._processId,
      state: this._recordedState,
      channelStates: {}
    }
    for (const [channel, messages] of this._recordedMessages.entries()) {
      snapshot.channelStates[`from_${channel.fromProcess._processId}`] = messages.map(
        msg => msg.content
      )
    }
    return snapshot
  }

  private _sendMarker() {
    for (const channel of this._outgoingChannels) {
      const marker: Message = { type: 'marker', content: '' }
      channel.send(marker)
      console.log(
        `Process ${this._processId} sent marker to Process ${channel.toProcess._processId}`
      )
    }
  }

  private _receiveMessage(channel: Channel, message: Message) {
    if (message.type === 'marker') {
      this._receiveMarker(channel)
    } else {
      if (this._isRecording && !this._receivedMarkers.has(channel)) {
        if (!this._recordedMessages.has(channel)) this._recordedMessages.set(channel, [])
        this._recordedMessages.get(channel)?.push(message)
      }
      // 处理普通消息
      console.log(
        `Process ${this._processId} received message from Process ${channel.fromProcess._processId}:`,
        message.content
      )
    }
  }

  private _receiveMarker(channel: Channel) {
    if (!this._isRecording) {
      // 第一次收到标记，开始记录
      this._isRecording = true
      this._recordedState = this._getState()
      console.log(`Process ${this._processId} records its state:`, this._recordedState)

      // 标记收到标记的通道
      this._receivedMarkers.add(channel)

      // 发送标记给所有输出通道
      this._sendMarker()

      // 开始记录其他输入通道的消息
      for (const ch of this._incomingChannels) {
        if (ch !== channel) {
          this._recordedMessages.set(ch, [])
        }
      }

      // 确保当前通道不用记录消息
      this._recordedMessages.set(channel, [])
    } else {
      // 已经开始记录，停止在该通道上的记录
      if (!this._receivedMarkers.has(channel)) {
        this._receivedMarkers.add(channel)
        console.log(
          `Process ${this._processId} stops recording on channel from Process ${channel.fromProcess._processId}`
        )
      }
    }
  }

  // 模拟获取当前进程状态
  private _getState(): any {
    // 这里可以根据实际需求定义状态
    return { counter: Math.floor(Math.random() * 100) }
  }
}

// 示例用法
const numProcesses = 3
const processes: Process[] = []

// 创建进程
for (let i = 0; i < numProcesses; i++) {
  processes.push(new Process(i))
}

// 创建通道，确保每个进程之间有单向通道
for (let i = 0; i < numProcesses; i++) {
  for (let j = 0; j < numProcesses; j++) {
    if (i !== j) {
      const channel = new Channel(processes[i], processes[j])
      processes[i].addOutgoingChannel(channel)
      processes[j].addIncomingChannel(channel)
    }
  }
}

// 启动进程
for (const process of processes) {
  process.run()
}

// 初始化进程状态并发送普通消息
for (const process of processes) {
  // 初始化状态
  // 这里的状态是随机生成的，可以根据实际需求修改
  // process.state = { counter: Math.floor(Math.random() * 100) };
}

// 发送一些普通消息
processes[0].sendMessage(processes[1], 'Hello from P0 to P1')
processes[1].sendMessage(processes[2], 'Hello from P1 to P2')
processes[2].sendMessage(processes[0], 'Hello from P2 to P0')

// 发起快照
setTimeout(() => {
  console.log('\n--- Initiating Snapshot at Process 0 ---')
  processes[0].startSnapshot()
}, 500)

// 记录快照结果
setTimeout(() => {
  console.log('\n--- Snapshot Results ---')
  const snapshots = processes.map(p => p.getSnapshot())
  console.log(JSON.stringify(snapshots, null, 2))
}, 2000)

export {}
