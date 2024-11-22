// # TypeScript 实现 Raft 共识算法

// ## 实现计划

// 要在 TypeScript 中实现 Raft 共识算法，需要按照以下步骤进行：

// 1. **定义基础类型和接口**
//    - 定义服务器角色（Leader、Follower、Candidate）。
//    - 定义日志条目结构。
//    - 定义 RPC 消息类型（RequestVote、AppendEntries 等）。

// 2. **实现服务器节点**
//    - 创建 `RaftNode` 类，包含当前任期、日志、状态、领导者信息等。
//    - 实现角色转换逻辑（Follower 到 Candidate，到 Leader）。

// 3. **实现选举过程**
//    - 处理选举超时，发起新的选举。
//    - 发送和处理 `RequestVote` RPC 消息。
//    - 处理投票请求并更新选票。

// 4. **实现日志复制**
//    - Leader 处理客户端请求，追加日志条目。
//    - Leader 发送 `AppendEntries` RPC 消息给 Followers。
//    - Followers 接收并追加日志，回复确认。

// 5. **实现心跳机制**
//    - Leader 定期发送空的 `AppendEntries` 消息作为心跳，防止 Followers 转变为 Candidate。

// 6. **处理故障和重新加入**
//    - 处理节点失败和恢复后的日志一致性。
//    - 使用快照（可选）来优化日志管理。

// 7. **实现网络通信**
//    - 使用 WebSocket 或其他通信方式在节点之间传递 RPC 消息。

// 8. **测试和验证**
//    - 编写测试用例，验证选举、日志复制、一致性等关键功能。

// ## TypeScript 代码实现

// 以下是 Raft 共识算法的简化实现示例：

// ### 1. 定义基础类型和接口

// // 定义服务器角色
enum Role {
  Follower,
  Candidate,
  Leader
}

// 日志条目结构
interface LogEntry {
  term: number
  command: any
}

// RequestVote RPC 消息
interface RequestVoteArgs {
  term: number
  candidateId: string
  lastLogIndex: number
  lastLogTerm: number
}

interface RequestVoteReply {
  term: number
  voteGranted: boolean
}

// AppendEntries RPC 消息
interface AppendEntriesArgs {
  term: number
  leaderId: string
  prevLogIndex: number
  prevLogTerm: number
  entries: LogEntry[]
  leaderCommit: number
}

interface AppendEntriesReply {
  term: number
  success: boolean
}

class RaftNode {
  id: string
  role: Role = Role.Follower
  currentTerm = 0
  votedFor: string | null = null
  log: LogEntry[] = []
  commitIndex = 0
  lastApplied = 0

  // Leader state
  nextIndex: Map<string, number> = new Map()
  matchIndex: Map<string, number> = new Map()

  // 状态
  peers: string[]
  state: any // 状态机

  // 定时器
  electionTimeout: NodeJS.Timeout | null = null
  heartbeatInterval: NodeJS.Timeout | null = null

  constructor(id: string, peers: string[]) {
    this.id = id
    this.peers = peers
    this.resetElectionTimeout()
  }

  // 重置选举超时
  resetElectionTimeout() {
    if (this.electionTimeout) clearTimeout(this.electionTimeout)
    const timeout = Math.random() * 150 + 150 // 150-300ms
    this.electionTimeout = setTimeout(() => this.startElection(), timeout)
  }

  // 开始选举
  startElection() {
    this.role = Role.Candidate
    this.currentTerm += 1
    this.votedFor = this.id
    let votesGranted = 1
    const args: RequestVoteArgs = {
      term: this.currentTerm,
      candidateId: this.id,
      lastLogIndex: this.log.length - 1,
      lastLogTerm: this.log[this.log.length - 1]?.term || 0
    }

    this.peers.forEach(peer => {
      this.sendRequestVote(peer, args).then(reply => {
        if (reply.voteGranted) votesGranted += 1
        if (votesGranted > this.peers.length / 2 && this.role === Role.Candidate) {
          this.becomeLeader()
        }
      })
    })

    this.resetElectionTimeout()
  }

  // 成为领导者
  becomeLeader() {
    this.role = Role.Leader
    this.peers.forEach(peer => {
      this.nextIndex.set(peer, this.log.length)
      this.matchIndex.set(peer, 0)
    })
    this.startHeartbeat()
  }

  // 启动心跳机制
  startHeartbeat() {
    if (this.heartbeatInterval) clearInterval(this.heartbeatInterval)
    this.heartbeatInterval = setInterval(() => this.sendHeartbeats(), 50)
  }

  // 发送心跳
  sendHeartbeats() {
    const args: AppendEntriesArgs = {
      term: this.currentTerm,
      leaderId: this.id,
      prevLogIndex: this.log.length - 1,
      prevLogTerm: this.log[this.log.length - 1]?.term || 0,
      entries: [],
      leaderCommit: this.commitIndex
    }

    this.peers.forEach(peer => {
      this.sendAppendEntries(peer, args).then(reply => {
        if (reply.term > this.currentTerm) {
          this.role = Role.Follower
          this.currentTerm = reply.term
          this.votedFor = null
          this.resetElectionTimeout()
          if (this.heartbeatInterval) clearInterval(this.heartbeatInterval)
        }
      })
    })
  }

  // 处理 RequestVote RPC
  async handleRequestVote(args: RequestVoteArgs): Promise<RequestVoteReply> {
    if (args.term < this.currentTerm) {
      return { term: this.currentTerm, voteGranted: false }
    }

    if (
      (this.votedFor === null || this.votedFor === args.candidateId) &&
      (args.lastLogTerm > this.log[this.log.length - 1]?.term ||
        (args.lastLogTerm === this.log[this.log.length - 1]?.term &&
          args.lastLogIndex >= this.log.length - 1))
    ) {
      this.votedFor = args.candidateId
      this.currentTerm = args.term
      this.resetElectionTimeout()
      return { term: this.currentTerm, voteGranted: true }
    }

    return { term: this.currentTerm, voteGranted: false }
  }

  // 处理 AppendEntries RPC
  async handleAppendEntries(args: AppendEntriesArgs): Promise<AppendEntriesReply> {
    if (args.term < this.currentTerm) {
      return { term: this.currentTerm, success: false }
    }

    this.role = Role.Follower
    this.currentTerm = args.term
    this.votedFor = args.leaderId
    this.resetElectionTimeout()

    if (
      args.prevLogIndex >= 0 &&
      (this.log.length <= args.prevLogIndex ||
        this.log[args.prevLogIndex].term !== args.prevLogTerm)
    ) {
      return { term: this.currentTerm, success: false }
    }

    // Append any new entries not already in the log
    let index = args.prevLogIndex + 1
    args.entries.forEach(entry => {
      if (this.log.length > index) {
        if (this.log[index].term !== entry.term) {
          this.log = this.log.slice(0, index)
          this.log.push(entry)
        }
      } else {
        this.log.push(entry)
      }
      index += 1
    })

    if (args.leaderCommit > this.commitIndex) {
      this.commitIndex = Math.min(args.leaderCommit, this.log.length - 1)
      this.applyLogs()
    }

    return { term: this.currentTerm, success: true }
  }

  // 应用日志到状态机
  applyLogs() {
    while (this.lastApplied < this.commitIndex) {
      this.lastApplied += 1
      const entry = this.log[this.lastApplied]
      // 应用到状态机
      this.state.apply(entry.command)
    }
  }

  // 发送 RequestVote RPC
  async sendRequestVote(peer: string, args: RequestVoteArgs): Promise<RequestVoteReply> {
    // 实现 RPC 调用，例如通过 WebSocket 或 HTTP
    // 这里使用伪代码
    return await rpcCall(peer, 'RequestVote', args)
  }

  // 发送 AppendEntries RPC
  async sendAppendEntries(peer: string, args: AppendEntriesArgs): Promise<AppendEntriesReply> {
    // 实现 RPC 调用，例如通过 WebSocket 或 HTTP
    // 这里使用伪代码
    return await rpcCall(peer, 'AppendEntries', args)
  }

  // 处理客户端请求（简化版）
  receiveClientCommand(command: any) {
    if (this.role !== Role.Leader) {
      // 重定向到领导者
      this.state.redirectToLeader()
      return
    }

    const entry: LogEntry = { term: this.currentTerm, command }
    this.log.push(entry)
    this.commitIndex += 1
    this.applyLogs()

    // 发送 AppendEntries 给 Followers
    const args: AppendEntriesArgs = {
      term: this.currentTerm,
      leaderId: this.id,
      prevLogIndex: this.log.length - 2,
      prevLogTerm: this.log.length >= 2 ? this.log[this.log.length - 2].term : 0,
      entries: [entry],
      leaderCommit: this.commitIndex
    }

    this.peers.forEach(peer => {
      this.sendAppendEntries(peer, args).then(reply => {
        if (reply.term > this.currentTerm) {
          this.role = Role.Follower
          this.currentTerm = reply.term
          this.votedFor = null
          this.resetElectionTimeout()
        }
      })
    })
  }
}
// 简化的 RPC 调用伪代码
async function rpcCall(peer: string, method: string, args: any): Promise<any> {
  // 使用 WebSocket、HTTP 或其他通信协议实现 RPC 调用
  // 这里只是一个伪实现
  const response = await fetch(`http://${peer}/raft/${method}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(args)
  })
  return await response.json()
}
class StateMachine {
  apply(command: any) {
    // 根据命令类型执行相应操作
    console.log(`Applying command: ${JSON.stringify(command)}`)
    // 具体业务逻辑
  }

  redirectToLeader() {
    // 实现重定向逻辑，例如返回领导者地址
    console.log('Redirecting to leader...')
  }
}
// 启动一个 Raft 节点
const nodeId = 'node1'
const peers = ['node2', 'node3']
const stateMachine = new StateMachine()

const raftNode = new RaftNode(nodeId, peers)
raftNode.state = stateMachine

// 示例：接收客户端命令
// raftNode.receiveClientCommand({ action: 'set', key: 'x', value: 10 });

export {}
