// 乐观更新 (Optimistic UI) —— 提升用户体验
// 业务场景：

// 点赞：点击红心，立马变红，然后再发请求。如果请求失败了，再变回灰色。
// 聊天发送：点击发送，消息立马出现在屏幕上（转圈圈），发送成功后圈圈消失。
// 核心痛点： 网络有延迟，如果等服务器返回再更新 UI，用户会觉得“卡”。

// 抽象实现 (Transaction Manager)： 我们需要一个机制来管理“临时状态”和“最终状态”。

export {}

// 简单的状态容器
class StateManager<T> {
  private state: T
  private listeners: ((s: T) => void)[] = []

  constructor(initial: T) {
    this.state = initial
  }

  getState() {
    return this.state
  }

  setState(newState: T) {
    this.state = newState
    this.listeners.forEach(fn => fn(this.state))
  }

  subscribe(fn: (s: T) => void) {
    this.listeners.push(fn)
  }

  // --- 核心：乐观更新 ---
  async performOptimisticUpdate(
    optimisticState: T, // 1. 预期的成功状态
    apiCall: () => Promise<void>, // 2. 实际的 API 请求
    rollbackState: T // 3. 失败后的回滚状态
  ) {
    // A. 立即应用乐观状态
    this.setState(optimisticState)

    try {
      // B. 发送请求
      await apiCall()
      // 请求成功，状态已经更新了，无需操作（或者用服务器返回的最新数据再次更新）
    } catch (error) {
      // C. 请求失败，回滚
      console.error('Update failed, rolling back...')
      this.setState(rollbackState)
    }
  }
}

// --- 业务实战：点赞 ---

const likeState = new StateManager({ liked: false, count: 100 })

// UI 绑定
likeState.subscribe(s => console.log(`[UI Render] Liked: ${s.liked}, Count: ${s.count}`))

// 用户点击点赞
const handleLike = () => {
  const current = likeState.getState()

  likeState.performOptimisticUpdate(
    // 1. 乐观值：立马变红，数字+1
    { liked: true, count: current.count + 1 },

    // 2. 慢速网络请求
    async () => {
      await new Promise((r, j) => setTimeout(Math.random() > 0.5 ? r : j, 1000)) // 50% 概率失败
    },

    // 3. 回滚值：变回原来的样子
    current
  )
}

handleLike()
