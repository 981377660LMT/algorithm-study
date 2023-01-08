/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable func-call-spacing */
/* eslint-disable no-spaced-func */

// 171. create callback on message channel
// This is a JavaScript coding problem from BFE.dev

interface SomePort {
  postMessage: (message: string) => void
  onmessage?: (message: string) => void
}

// 现在SomeChannel传递信息存在随机延迟,这意味着
// port发送消息的顺序可能与接收消息的顺序不同。
// !需要使用 回调函数(事件监听)+map(消息id=>回调) 异步来解决channel的消息顺序问题
declare class SomeChannel {
  port1: SomePort
  port2: SomePort
}

/**
 *
  Implement a BetterChannel class that enables communication between ports
  `with callback` and reply handle, but on top of SomeChannel.

  ```js
  const {port1, port2} = new BetterChannel()

  port2.onmessage = (message, reply) => {
    if (message === 'ping?') {
      reply('pong!')
    }
    if (message === 'pong?') {
      reply('ping!')
    }
  }

  port1.postMessage('ping?', (data) => {
    console.log(data) // 'pong!'
   })
   ```
 */
class BetterChannel {
  readonly port1: WrappedSomePort
  readonly port2: WrappedSomePort
  private _id = 0 // can be replaced by uuid
  private readonly _cbRecord = new Map<number, Callback>()

  constructor() {
    const badChannel = new SomeChannel()
    this.port1 = {} as WrappedSomePort
    this.port2 = {} as WrappedSomePort
    this.port1.postMessage = this._createPostMessage(badChannel.port1, badChannel.port2, this.port2)
    this.port2.postMessage = this._createPostMessage(badChannel.port2, badChannel.port1, this.port1)
  }

  // 调用对方的onmessage,返回自己的postMessage
  // TODO: 为什么要这么做?
  private _createPostMessage(
    sender: SomePort,
    receiver: SomePort,
    peer: WrappedSomePort
  ): WrappedSomePort['postMessage'] {
    receiver.onmessage = (message: string) => {
      if (!peer.onmessage) return
      const { id, message: msg } = JSON.parse(message)
      peer.onmessage(msg, (response: unknown) => {
        if (!this._cbRecord.has(id)) return
        const cb = this._cbRecord.get(id)!
        cb(response)
        this._cbRecord.delete(id)
      })
    }

    return (message: string, cb: Callback) => {
      const id = this._id++
      this._cbRecord.set(id, cb)
      sender.postMessage(JSON.stringify({ id, message }))
    }
  }
}

type Callback = (response: unknown) => void

interface WrappedSomePort {
  postMessage: (message: string, cb: Callback) => void
  onmessage?: (message: string, cb: Callback) => void
}

if (require.main === module) {
  const { port1, port2 } = new BetterChannel()
  port1.postMessage('ping?', data => {
    console.log(data) // 'pong!'
  })

  port2.onmessage = (message, reply) => {
    if (message === 'ping?') {
      reply('pong!')
    }
    if (message === 'pong?') {
      reply('ping!')
    }
  }
}

export {}
