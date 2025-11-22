// 熔断器 (Circuit Breaker) —— 服务的自我保护
// 业务场景：
// 非关键服务挂了：你的页面里有一个“猜你喜欢”的侧边栏，调用的推荐接口非常慢或者频繁 500。你不想因为这个次要功能导致整个页面卡死（Loading 转个不停）。
// 微服务/BFF 调用：Node.js 中间层调用下游 Java 服务，下游挂了，Node 层连接池被占满，导致整个 Node 服务不可用。
// !核心痛点： 当某个依赖服务出现故障时，持续的重试只会加重负担（雪崩效应），且拖累主流程。我们需要“快速失败”。
// !抽象实现： 维护三种状态：关闭 (Closed)（正常）、打开 (Open)（熔断，直接报错）、半开 (Half-Open)（尝试恢复）。

enum CircuitState {
  CLOSED, // 正常
  OPEN, // 熔断
  HALF_OPEN // 尝试恢复
}

export class CircuitBreaker {
  private state = CircuitState.CLOSED
  private failureCount = 0
  private lastFailureTime = 0

  constructor(
    private threshold = 3, // 失败多少次触发熔断
    private recoveryTimeout = 5000 // 熔断多久后尝试恢复
  ) {}

  async call<T>(fn: () => Promise<T>): Promise<T> {
    // 1. 检查是否熔断
    if (this.state === CircuitState.OPEN) {
      if (Date.now() - this.lastFailureTime > this.recoveryTimeout) {
        this.state = CircuitState.HALF_OPEN // 时间到了，试一试
      } else {
        throw new Error('Circuit is OPEN: Fast fail')
      }
    }

    try {
      // 2. 执行函数
      const result = await fn()

      // 3. 成功回调
      if (this.state === CircuitState.HALF_OPEN) {
        this.state = CircuitState.CLOSED // 恢复成功，关闭熔断器
        this.failureCount = 0
      }
      return result
    } catch (err) {
      // 4. 失败处理
      this.failureCount++
      this.lastFailureTime = Date.now()

      if (this.failureCount >= this.threshold) {
        this.state = CircuitState.OPEN // 失败次数超标，开启熔断
        console.warn('Circuit Breaker Tripped!')
      }
      throw err
    }
  }
}

// --- 业务实战：调用推荐接口 ---

const breaker = new CircuitBreaker(3, 2000) // 3次失败熔断，2秒后重试

async function fetchRecommendations() {
  // 模拟不稳定的 API
  if (Math.random() > 0.5) throw new Error('API Timeout')
  return ['Item A', 'Item B']
}

// 页面加载逻辑
async function loadPage() {
  try {
    // 即使 fetchRecommendations 挂了，breaker 会快速抛错，不会卡住 await
    const data = await breaker.call(fetchRecommendations)
    console.log('Render:', data)
  } catch (e) {
    console.log('Render: Default Recommendations (Fallback)')
  }
}
