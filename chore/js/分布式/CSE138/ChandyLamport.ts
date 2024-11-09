// chandy-lamport.ts

class Message {
  content: string
  isMarker: boolean

  constructor(content: string, isMarker = false) {
    this.content = content
    this.isMarker = isMarker
  }
}

class Process {
  private readonly _processId: number
  private readonly _numProcesses: number
  private readonly _channels: Message[][]
  private readonly _markerReceived: boolean[]
  private _state: string | undefined

  constructor(processId: number, numProcesses: number) {
    this._processId = processId
    this._numProcesses = numProcesses
    this._channels = Array(numProcesses)
    for (let i = 0; i < numProcesses; i++) this._channels[i] = []
    this._markerReceived = Array(numProcesses).fill(false)
    this._state = undefined
  }

  run() {
    setInterval(() => {
      for (let i = 0; i < this._numProcesses; i++) {
        if (i !== this._processId) {
          const message = this._receiveMessage(processes[i])
          if (message) {
            if (message.isMarker) {
              this._receiveMarker(processes[i])
            } else {
              console.log(
                `Process ${this._processId} received message from Process ${i}: ${message.content}`
              )
            }
          }
        }
      }
    }, 100)
  }

  sendMessage(targetProcess: Process, message: Message) {
    console.log(
      `Process ${this._processId} sending message to Process ${targetProcess._processId}: ${message.content}`
    )
    this._channels[targetProcess._processId].push(message)
  }

  startSnapshot() {
    this._recordState()
    for (let i = 0; i < this._numProcesses; i++) {
      if (i !== this._processId) {
        this.sendMessage(processes[i], new Message('Marker', true))
      }
    }
  }

  private _receiveMessage(sourceProcess: Process): Message | undefined {
    return this._channels[sourceProcess._processId].shift()
  }

  private _receiveMarker(sourceProcess: Process) {
    if (!this._markerReceived[sourceProcess._processId]) {
      this._markerReceived[sourceProcess._processId] = true
      this._recordState()
      for (let i = 0; i < this._numProcesses; i++) {
        if (i !== this._processId) {
          this.sendMessage(processes[i], new Message('Marker', true))
        }
      }
    }
    this._recordChannelState(sourceProcess)
  }

  private _recordState() {
    this._state = `State of Process ${this._processId} at time ${Date.now()}`
    console.log(this._state)
  }

  private _recordChannelState(sourceProcess: Process) {
    const messages: string[] = []
    const channel = this._channels[sourceProcess._processId]
    if (channel) {
      while (channel.length > 0) {
        messages.push(channel.shift()!.content)
      }
    }
    console.log(
      `Channel state from Process ${sourceProcess._processId} to Process ${this._processId}: ${messages}`
    )
  }
}

const processes: Process[] = []

function main() {
  const numProcesses = 3
  for (let i = 0; i < numProcesses; i++) {
    processes.push(new Process(i, numProcesses))
  }

  // Start processes
  processes.forEach(process => process.run())

  // Simulate message passing
  setTimeout(() => {
    processes[0].sendMessage(processes[1], new Message('Hello from P0 to P1'))
    processes[1].sendMessage(processes[2], new Message('Hello from P1 to P2'))
    processes[2].sendMessage(processes[0], new Message('Hello from P2 to P0'))
  }, 1000)

  // Start snapshot
  setTimeout(() => {
    processes[0].startSnapshot()
  }, 2000)
}

main()
