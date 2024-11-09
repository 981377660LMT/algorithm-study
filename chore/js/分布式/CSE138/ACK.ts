interface Message {
  senderId: string
  sequenceNumber: number
  content: string
}

interface Ack {
  senderId: string
  sequenceNumber: number
}

class FIFOChannel {
  private senderToReceiver: Message[] = []
  private receiverToSender: Ack[] = []

  transmitMessage(message: Message) {
    console.log(`Channel: Transmitting message ${JSON.stringify(message)}`)
    this.senderToReceiver.push(message)
  }

  transmitAck(ack: Ack) {
    console.log(`Channel: Transmitting ACK ${JSON.stringify(ack)}`)
    this.receiverToSender.push(ack)
  }

  receiveMessage(): Message | undefined {
    return this.senderToReceiver.shift()
  }

  receiveAck(): Ack | undefined {
    return this.receiverToSender.shift()
  }
}

class Sender {
  private sequenceNumber: number = 0
  private expectedAck: number = -1
  private ackReceived: boolean = false

  constructor(
    private senderId: string,
    private channel: FIFOChannel,
    private timeout: number = 2000 // 超时时间，单位毫秒
  ) {
    this.startAckListener()
  }

  send(content: string) {
    const message: Message = {
      senderId: this.senderId,
      sequenceNumber: this.sequenceNumber,
      content
    }
    this.expectedAck = this.sequenceNumber
    this.ackReceived = false

    this.channel.transmitMessage(message)
    console.log(`Sender ${this.senderId}: Sent ${JSON.stringify(message)}`)

    // 等待ACK
    setTimeout(() => {
      if (!this.ackReceived) {
        console.log(
          `Sender ${this.senderId}: ACK timeout for seq ${this.sequenceNumber}, resending...`
        )
        this.channel.transmitMessage(message)
        console.log(`Sender ${this.senderId}: Resent ${JSON.stringify(message)}`)
      }
    }, this.timeout)

    this.sequenceNumber++
  }

  private startAckListener() {
    setInterval(() => {
      const ack = this.channel.receiveAck()
      if (ack && ack.senderId === this.senderId && ack.sequenceNumber === this.expectedAck) {
        this.ackReceived = true
        console.log(`Sender ${this.senderId}: ACK received for seq ${ack.sequenceNumber}`)
      }
    }, 500)
  }
}

class Receiver {
  private expectedSeq: number = 0

  constructor(private channel: FIFOChannel) {
    this.startMessageListener()
  }

  private startMessageListener() {
    setInterval(() => {
      const message = this.channel.receiveMessage()
      if (message) {
        console.log(`Receiver: Received ${JSON.stringify(message)}`)
        if (message.sequenceNumber === this.expectedSeq) {
          this.deliver(message)
          const ack: Ack = {
            senderId: message.senderId,
            sequenceNumber: message.sequenceNumber
          }
          this.channel.transmitAck(ack)
          console.log(`Receiver: Sent ACK for seq ${ack.sequenceNumber}`)
          this.expectedSeq += 1
        } else if (message.sequenceNumber < this.expectedSeq) {
          // 重复或过时的消息，重新发送ACK
          const ack: Ack = {
            senderId: message.senderId,
            sequenceNumber: message.sequenceNumber
          }
          this.channel.transmitAck(ack)
          console.log(`Receiver: Resent ACK for seq ${ack.sequenceNumber}`)
        } else {
          // 未来的消息，忽略或处理
          console.log(
            `Receiver: Out-of-order message ${JSON.stringify(message)}, expected seq ${this.expectedSeq}. Ignoring.`
          )
        }
      }
    }, 500)
  }

  private deliver(message: Message) {
    console.log(`Receiver: Delivered ${JSON.stringify(message)} to application.`)
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  // 创建频道、接收者和发送者
  const channel = new FIFOChannel()
  const receiver = new Receiver(channel)
  const sender = new Sender('P1', channel, 3000) // 超时时间设置为3秒

  // 发送有序消息
  sender.send('Message 1')
  sender.send('Message 2')
  sender.send('Message 3')

  console.log('\n--- 模拟ACK丢失导致重传 ---\n')

  // 修改FIFOChannel以模拟ACK丢失
  class UnreliableFIFOChannel extends FIFOChannel {
    override transmitAck(ack: Ack) {
      if (ack.sequenceNumber === 1) {
        // 模拟ACK丢失
        console.log(`Channel: ACK ${JSON.stringify(ack)} lost!`)
        return
      }
      super.transmitAck(ack)
    }
  }

  // 使用不可靠的频道
  const unreliableChannel = new UnreliableFIFOChannel()
  const unreliableReceiver = new Receiver(unreliableChannel)
  const unreliableSender = new Sender('P1', unreliableChannel, 3000)

  unreliableSender.send('Message 4') // ACK将丢失，触发重传
  unreliableSender.send('Message 5')

  // 等待一段时间以观察重传效果
  setTimeout(() => {
    console.log('结束模拟.')
    process.exit(0)
  }, 10000)
}
