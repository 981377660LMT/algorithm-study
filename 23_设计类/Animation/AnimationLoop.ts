/**
 * 动画状态的泛型。
 */
type TState = any

/**
 * 状态更新函数。
 * @param state - 当前状态。
 * @param deltaTime - 自上一帧以来经过的时间（以毫秒为单位）。
 * @returns 返回下一帧的新状态。
 */
type UpdateCallback<S extends TState> = (state: S, deltaTime: number) => S

/**
 * 渲染函数。
 * @param state - 需要被渲染的当前状态。
 */
type RenderCallback<S extends TState> = (state: S) => void

/**
 * 调度器需要调用的回调函数类型。
 */
type TickCallback = () => void

/**
 * 调度器接口。
 * 它的职责是根据自己的策略（RAF, setTimeout, 手动等）
 * 来重复调用一个 tick 回调函数。
 */
interface IScheduler {
  /**
   * 启动调度器。
   * @param tick - 每一帧需要执行的回调函数。
   */
  start(tick: TickCallback): void

  /**
   * 停止调度器。
   */
  stop(): void
}

class RequestAnimationFrameScheduler implements IScheduler {
  private animationFrameId: number | null = null

  public start(tick: TickCallback): void {
    const loop = () => {
      tick()
      this.animationFrameId = requestAnimationFrame(loop)
    }
    // 立即开始第一次循环
    this.animationFrameId = requestAnimationFrame(loop)
  }

  public stop(): void {
    if (this.animationFrameId) {
      cancelAnimationFrame(this.animationFrameId)
      this.animationFrameId = null
    }
  }
}

/** 固定帧率. */
class SetTimeoutScheduler implements IScheduler {
  private timeoutId: number | null = null
  private readonly fps: number

  constructor(fps: number = 60) {
    this.fps = fps
  }

  public start(tick: TickCallback): void {
    const interval = 1000 / this.fps
    const loop = () => {
      tick()
      this.timeoutId = window.setTimeout(loop, interval)
    }
    // 立即开始第一次循环
    this.timeoutId = window.setTimeout(loop, interval)
  }

  public stop(): void {
    if (this.timeoutId) {
      clearTimeout(this.timeoutId)
      this.timeoutId = null
    }
  }
}

export class AnimationLoop<S extends TState> {
  private state: S
  private update: UpdateCallback<S>
  private render: RenderCallback<S>
  private scheduler: IScheduler

  private lastTime: number = 0
  private isRunning: boolean = false

  /**
   * 创建一个动画循环实例。
   * @param initialState - 动画的初始状态。
   * @param update - 每一帧用于计算新状态的回调函数。
   * @param render - 每一帧用于渲染状态的回调函数。
   * @param scheduler - 调度器实例，用于控制动画的帧率。默认为 RequestAnimationFrameScheduler。
   */
  constructor(
    initialState: S,
    update: UpdateCallback<S>,
    render: RenderCallback<S>,
    scheduler: IScheduler = new RequestAnimationFrameScheduler()
  ) {
    this.state = initialState
    this.update = update
    this.render = render
    this.scheduler = scheduler
  }

  public start(): void {
    if (this.isRunning) {
      return
    }
    this.isRunning = true
    this.lastTime = performance.now()
    // 将 tick 方法传递给调度器，让它来控制循环
    this.scheduler.start(this.tick)
    console.log('Animation started.')
  }

  public stop(): void {
    if (!this.isRunning) {
      return
    }
    this.isRunning = false
    // 同样，让调度器来停止循环
    this.scheduler.stop()
    console.log('Animation stopped.')
  }

  public getIsRunning(): boolean {
    return this.isRunning
  }

  /**
   * 将由调度器调用。
   */
  private tick = (): void => {
    if (!this.isRunning) {
      return
    }

    const currentTime = performance.now()
    const deltaTime = currentTime - this.lastTime
    this.lastTime = currentTime

    this.state = this.update(this.state, deltaTime)
    this.render(this.state)
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  // 1. 定义业务状态
  interface CircleState {
    x: number
    y: number
    radius: number
    vx: number // x轴速度 (像素/毫秒)
    vy: number // y轴速度 (像素/毫秒)
  }

  // 2. 获取 Canvas 上下文
  const canvas = document.createElement('canvas')
  canvas.width = 400
  canvas.height = 300
  document.body.appendChild(canvas)
  const ctx = canvas.getContext('2d')!

  // 3. 定义业务逻辑的回调函数
  const initialState: CircleState = {
    x: 50,
    y: 50,
    radius: 20,
    vx: 0.1,
    vy: 0.1
  }

  // 状态更新逻辑
  const updateCircle: UpdateCallback<CircleState> = (state, deltaTime) => {
    // 创建一个新状态，而不是修改旧状态（遵循不可变性原则）
    const newState = { ...state }

    newState.x += newState.vx * deltaTime
    newState.y += newState.vy * deltaTime

    // 边界碰撞检测
    if (newState.x < state.radius || newState.x > canvas.width - state.radius) {
      newState.vx *= -1
    }
    if (newState.y < state.radius || newState.y > canvas.height - state.radius) {
      newState.vy *= -1
    }

    return newState
  }

  // 渲染逻辑
  const renderCircle: RenderCallback<CircleState> = state => {
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    ctx.beginPath()
    ctx.arc(state.x, state.y, state.radius, 0, Math.PI * 2)
    ctx.fillStyle = '#28b'
    ctx.fill()
  }

  // 4. 创建并启动动画
  const fixedFpsScheduler = new SetTimeoutScheduler(30)

  // 将调度器实例传入 AnimationLoop 的构造函数
  const choppyAnimation = new AnimationLoop(
    initialState,
    updateCircle,
    renderCircle,
    fixedFpsScheduler
  )
  choppyAnimation.start()
}
