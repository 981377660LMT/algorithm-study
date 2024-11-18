/* eslint-disable no-case-declarations */
// TypeScript 实现 Paxos 算法

type Message<T> =
  | PrepareMessage
  | PromiseMessage<T>
  | AcceptRequestMessage<T>
  | AcceptMessage<T>
  | LearnMessage<T>

enum MessageType {
  PREPARE = 0,
  PROMISE,
  ACCEPT_REQUEST,
  ACCEPTED,
  LEARN
}

interface PrepareMessage {
  from: number
  to: number
  type: MessageType.PREPARE

  proposalNumber: number
}

interface PromiseMessage<T> {
  from: number
  to: number
  type: MessageType.PROMISE

  proposalNumber: number
  acceptedProposal: number
  acceptedValue: T | undefined
}

interface AcceptRequestMessage<T> {
  from: number
  to: number
  type: MessageType.ACCEPT_REQUEST

  proposalNumber: number
  value: T | undefined
}

interface AcceptMessage<T> {
  from: number
  to: number
  type: MessageType.ACCEPTED

  proposalNumber: number
  value: T | undefined
}

interface LearnMessage<T> {
  from: number
  to: number
  type: MessageType.LEARN

  proposalNumber: number
  value: T | undefined
}

interface Node<T> {
  id: number
  receive(message: Message<T>): void
}

class Proposer<T> implements Node<T> {
  readonly id: number

  private readonly _network: Network<T>
  private readonly _majority: number // 多数阈值

  private _proposalNumber = 0 // 提案编号
  private _proposalValue: T | undefined = undefined // 提案值
  private readonly _promisesReceived: Map<number, PromiseMessage<T>> = new Map() // 收到的承诺

  constructor(id: number, network: Network<T>, majority: number) {
    this.id = id
    this._network = network
    this._majority = majority
  }

  propose(value: T) {
    this._proposalNumber++
    this._proposalValue = value
    const prepareMessage: Message<T> = {
      from: this.id,
      to: -1,
      type: MessageType.PREPARE,
      proposalNumber: this._proposalNumber
    }
    this._network.broadcastMessage(prepareMessage)
  }

  receive(message: Message<T>) {
    switch (message.type) {
      case MessageType.PROMISE:
        this._handlePromise(message)
        break
      case MessageType.ACCEPTED:
        this._handleAccepted(message)
        break
      default:
        break
    }
  }

  private _handlePromise(message: PromiseMessage<T>) {
    const n = message.proposalNumber
    if (n === this._proposalNumber) {
      this._promisesReceived.set(message.from, message)
      if (this._promisesReceived.size >= this._majority) {
        let highestAcceptedProposal = -1
        for (const promise of this._promisesReceived.values()) {
          if (promise.acceptedProposal > highestAcceptedProposal) {
            highestAcceptedProposal = promise.acceptedProposal
            if (promise.acceptedValue !== undefined) {
              this._proposalValue = promise.acceptedValue
            }
          }
        }

        const acceptMessage: Message<T> = {
          from: this.id,
          to: -1,
          type: MessageType.ACCEPT_REQUEST,
          proposalNumber: this._proposalNumber,
          value: this._proposalValue
        }
        this._network.broadcastMessage(acceptMessage)
      }
    }
  }

  // eslint-disable-next-line class-methods-use-this
  private _handleAccepted(message: AcceptMessage<T>) {}
}

class Acceptor<T> implements Node<T> {
  readonly id: number
  private readonly _network: Network<T>

  private _promisedProposal = 0 // 承诺过的最大提案编号
  private _acceptedProposal = 0 // 已接受的最大提案编号
  private _acceptedValue: T | undefined = undefined // 已接受的最大编号对应的值

  constructor(id: number, network: Network<T>) {
    this.id = id
    this._network = network
  }

  receive(message: Message<T>) {
    switch (message.type) {
      case MessageType.PREPARE:
        this._handlePrepare(message)
        break
      case MessageType.ACCEPT_REQUEST:
        this._handleAcceptRequest(message)
        break
      default:
        break
    }
  }

  private _handlePrepare(message: PrepareMessage) {
    const n = message.proposalNumber
    if (n > this._promisedProposal) {
      this._promisedProposal = n
      const reply: Message<T> = {
        from: this.id,
        to: message.from,
        type: MessageType.PROMISE,
        proposalNumber: n,
        acceptedProposal: this._acceptedProposal,
        acceptedValue: this._acceptedValue
      }
      this._network.sendMessage(reply)
    }
  }

  private _handleAcceptRequest(message: AcceptRequestMessage<T>) {
    const n = message.proposalNumber
    if (n >= this._promisedProposal) {
      this._promisedProposal = n
      this._acceptedProposal = n
      this._acceptedValue = message.value
      const reply: Message<T> = {
        from: this.id,
        to: message.from,
        type: MessageType.ACCEPTED,
        proposalNumber: n,
        value: this._acceptedValue
      }
      this._network.sendMessage(reply)
      this._notifyLearners(n, this._acceptedValue)
    }
  }

  private _notifyLearners(n: number, value: T | undefined) {
    const learnMessage: Message<T> = {
      value,
      proposalNumber: n,
      from: this.id,
      to: -1,
      type: MessageType.LEARN
    }
    this._network.broadcastMessage(learnMessage)
  }
}

class Learner<T> implements Node<T> {
  readonly id: number
  private readonly _learnedValues: (T | undefined)[] = []

  constructor(id: number) {
    this.id = id
  }

  receive(message: Message<T>) {
    if (message.type === MessageType.LEARN) {
      this._learnedValues.push(message.value)
      console.log(`学习者 ${this.id} 学到了值: ${message.value}`)
    }
  }
}

class Network<T> {
  private readonly _nodes: Map<number, Node<T>> = new Map()

  registerNode(node: Node<T>) {
    this._nodes.set(node.id, node)
  }

  sendMessage(message: Message<T>) {
    if (message.to === -1) {
      // 广播消息
      for (const node of this._nodes.values()) {
        if (node.id !== message.from) {
          node.receive(message)
        }
      }
    } else {
      const node = this._nodes.get(message.to)
      if (node) {
        node.receive(message)
      }
    }
  }

  broadcastMessage(message: Message<T>) {
    message.to = -1
    this.sendMessage(message)
  }
}

// 示例使用
const network = new Network()

// 创建接受者（Acceptor）
for (let i = 1; i <= 3; i++) {
  const acceptor = new Acceptor(i, network)
  network.registerNode(acceptor)
}

// 创建学习者（Learner）
for (let i = 4; i <= 5; i++) {
  const learner = new Learner(i)
  network.registerNode(learner)
}

// 创建提议者（Proposer）
const proposer = new Proposer(-1, network, 2) // 多数为2
network.registerNode(proposer)

const proposer2 = new Proposer(-2, network, 2) // 多数为2
network.registerNode(proposer2)

// 发起提案
proposer.propose('值_A')
proposer2.propose('值_B')

export {}
