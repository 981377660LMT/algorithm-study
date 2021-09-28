import { MinHeap } from '../2_queue/minheap'

class FileSharing {
  private chunkToUser: Map<number, Set<number>>
  private userToChunk: Map<number, Set<number>>
  private releasedUserId: MinHeap
  private nextUserId: number

  constructor(private m: number) {
    this.chunkToUser = new Map()
    this.userToChunk = new Map()
    this.releasedUserId = new MinHeap()
    this.nextUserId = 1
    for (let i = 1; i <= m; i++) {
      this.chunkToUser.set(i, new Set())
    }
  }

  // 系统应为其注册一个独有的 ID。这个独有的 ID 应当被相应的用户使用一次，
  join(ownedChunks: number[]): number {
    let userId = this.nextUserId
    if (this.releasedUserId.size) userId = this.releasedUserId.shift()!
    else this.nextUserId++

    this.userToChunk.set(userId, new Set(ownedChunks))
    ownedChunks.forEach(chunk => {
      this.chunkToUser.get(chunk)?.add(userId)
    })

    return userId
  }

  // 当用户离开系统时，其 ID 应可以被（后续新注册的用户）再次使用
  leave(userID: number): void {
    this.releasedUserId.push(userID)
    const chunks = this.userToChunk.get(userID) || []
    chunks.forEach(chunk => {
      this.chunkToUser.get(chunk)?.delete(userID)
    })
  }

  // 系统应当返回拥有这个文件块的所有用户的 ID,按升序排列。如果用户收到 ID 的非空列表，就表示成功接收到请求的文件块
  // 当某个用户请求了一个文件块时，他也就拥有了那个文件块
  request(userID: number, chunkID: number): number[] {
    const requestResult = this.chunkToUser.get(chunkID)
    if (!requestResult) return []
    const res = [...requestResult].sort((a, b) => a - b)
    if (requestResult.size && !requestResult.has(userID)) {
      this.userToChunk.get(userID)?.add(chunkID)
      this.chunkToUser.get(chunkID)?.add(userID)
    }
    return res
  }
}

const u = new FileSharing(4)
u.join([1, 2])
u.join([2, 3])
u.join([4])
console.log(u)
console.log(u.request(1, 3))
export {}

// 我们需要使用一套文件分享系统来分享一个非常大的文件，该文件由 m 个从 1 到 m 编号的文件块组成
// 1、request的时候要判断返回列表是否为空，如果为空的话，表示没有人拥有文件，那么就返回为空；
// 2、如果不为空，那么这次请求以后，这个人也拥有了文件，就把这个人加到文件拥有者列表。

// 总结：
// 1.使用优先队列维护回收的id
// 如果非空 则优先去回收的id里取
// 2.由于返回的结果要排序 应当使用treeSet而非Set
