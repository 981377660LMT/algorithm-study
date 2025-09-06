/**
 * 动画状态的泛型。
 */
type TState = any

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

/**
 * 代表一个可动画化对象的接口。
 * 它封装了状态以及更新和渲染该状态的逻辑。
 */
interface IAnimatable<S extends TState> {
  state: S
  update(deltaTime: number): void
  render(): void
}

/**
 * 定义了单帧执行策略的接口。
 */
interface IFrameRunner<S extends TState> {
  /**
   * 执行一帧的逻辑。
   * @param animatable - 需要执行动画的对象。
   * @param deltaTime - 时间增量。
   */
  runFrame(animatable: IAnimatable<S>, deltaTime: number): void
}

class SimpleFrameRunner<S extends TState> implements IFrameRunner<S> {
  public runFrame(animatable: IAnimatable<S>, deltaTime: number): void {
    animatable.update(deltaTime)
    animatable.render()
  }
}

class FixedTimestepRunner<S extends TState> implements IFrameRunner<S> {
  private accumulator: number = 0
  private readonly fixedDeltaTime: number

  constructor(updatesPerSecond: number = 60) {
    this.fixedDeltaTime = 1000 / updatesPerSecond
  }

  public runFrame(animatable: IAnimatable<S>, deltaTime: number): void {
    // 防止因浏览器标签页切换等原因导致 deltaTime 过大
    if (deltaTime > 250) {
      deltaTime = 250
    }
    this.accumulator += deltaTime

    while (this.accumulator >= this.fixedDeltaTime) {
      animatable.update(this.fixedDeltaTime)
      this.accumulator -= this.fixedDeltaTime
    }

    animatable.render()
  }
}

/**
 * `AnimationRunner` 彻底变成了一个通用的“协调器”，它对动画的具体内容、执行时机、执行策略一无所知.
 */
export class AnimationRunner<S extends TState> {
  private animatable: IAnimatable<S>
  private scheduler: IScheduler
  private frameRunner: IFrameRunner<S>

  private lastTime: number = 0
  private isRunning: boolean = false

  constructor(
    animatable: IAnimatable<S>,
    scheduler: IScheduler = new RequestAnimationFrameScheduler(),
    frameRunner: IFrameRunner<S> = new SimpleFrameRunner<S>()
  ) {
    this.animatable = animatable
    this.scheduler = scheduler
    this.frameRunner = frameRunner
  }

  public start(): void {
    if (this.isRunning) return
    this.isRunning = true
    this.lastTime = performance.now()
    this.scheduler.start(this.tick)
    console.log('Animation started.')
  }

  public stop(): void {
    if (!this.isRunning) return
    this.isRunning = false
    this.scheduler.stop()
    console.log('Animation stopped.')
  }

  private tick = (): void => {
    if (!this.isRunning) return

    const currentTime = performance.now()
    const deltaTime = currentTime - this.lastTime
    this.lastTime = currentTime

    // 将所有工作委托给 FrameRunner
    this.frameRunner.runFrame(this.animatable, deltaTime)
  }
}

if (require.main === module) {
  interface CircleState {
    x: number
    y: number
    radius: number
    vx: number
    vy: number
  }

  class BouncingCircle implements IAnimatable<CircleState> {
    public state: CircleState
    private ctx: CanvasRenderingContext2D

    constructor(initialState: CircleState, context: CanvasRenderingContext2D) {
      this.state = initialState
      this.ctx = context
    }

    update(deltaTime: number): void {
      this.state.x += this.state.vx * deltaTime
      this.state.y += this.state.vy * deltaTime

      const canvas = this.ctx.canvas
      if (this.state.x < this.state.radius || this.state.x > canvas.width - this.state.radius) {
        this.state.vx *= -1
      }
      if (this.state.y < this.state.radius || this.state.y > canvas.height - this.state.radius) {
        this.state.vy *= -1
      }
    }

    render(): void {
      const { width, height } = this.ctx.canvas
      this.ctx.clearRect(0, 0, width, height)
      this.ctx.beginPath()
      this.ctx.arc(this.state.x, this.state.y, this.state.radius, 0, Math.PI * 2)
      this.ctx.fillStyle = '#28b'
      this.ctx.fill()
    }
  }

  const canvas = document.createElement('canvas')
  document.body.appendChild(canvas)
  const ctx = canvas.getContext('2d')!
  const initialState: CircleState = { x: 50, y: 50, radius: 20, vx: 0.1, vy: 0.1 }

  // 创建 Animatable 实例
  const myCircle = new BouncingCircle(initialState, ctx)

  // 组合1: 默认行为 (RAF + 简单策略)
  const runner1 = new AnimationRunner(myCircle)
  runner1.start()

  // 组合2: 固定 30fps + 稳定逻辑更新策略
  const runner2 = new AnimationRunner(
    myCircle,
    new SetTimeoutScheduler(30),
    new FixedTimestepRunner(60) // 逻辑更新频率为 60/s
  )
  // runner2.start();
}
