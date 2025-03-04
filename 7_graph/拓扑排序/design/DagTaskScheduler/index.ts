// DAGScheduler
// TODO: test

/**
 * 任务接口，定义了可以被DAG调度器执行的任务的基本结构.
 * @template C - 上下文类型，用于在任务执行时传递共享数据.
 */
interface ITask<C> {
  /**
   * 任务的唯一标识符，用于在DAG中引用和识别任务.
   * 必须在整个DAG中保持唯一，否则添加时会引发错误.
   */
  readonly id: string

  /**
   * 当前任务依赖的其他任务ID数组.
   * @example ['task1', 'task2']，表示当前任务需要等待task1和task2执行完成后才能执行.
   */
  readonly deps: string[]

  /**
   * 当任务被触发执行时调用此方法.
   * 任务的主要业务逻辑应该在这里实现.
   */
  onTrigger(context: C): void | Promise<void>

  /**
   * 当任务需要重置状态时调用此方法.
   * 通常在其依赖的任务被重新触发后，此任务也需要重置状态.
   */
  onReset(context: C): void | Promise<void>

  /**
   * 当任务执行出错时的错误处理方法.
   * 用于实现自定义错误处理逻辑，如记录日志或清理资源等.
   */
  onError(context: C, error: Error): void | Promise<void>
}

/**
 * 表示任务在DAG调度器中的执行状态.
 *
 * - 'idle'：初始状态或已重置状态.表示任务尚未开始执行或已被重置为初始状态，
 *   等待被触发.任务初始创建或调用onReset后会处于此状态.
 *
 * - 'pending'：执行中状态.表示任务已被触发并正在执行，尚未完成.
 *   任务执行过程中（onTrigger或onReset调用期间）处于此状态.
 *
 * - 'completed'：已完成状态.表示任务已成功完成执行.
 *   只有当任务的onTrigger方法成功完成后才会进入此状态.
 *
 * - 'errored'：错误状态.表示任务执行过程中发生了错误.
 *   当onTrigger或onReset方法抛出异常时会进入此状态.
 */
type TaskStatus = 'idle' | 'pending' | 'completed' | 'errored'

class TaskNode<C> {
  readonly deps: Set<string>
  readonly children: Set<string>

  private readonly _task: ITask<C>

  private _status: TaskStatus

  constructor(task: ITask<C>) {
    this.deps = new Set(task.deps)
    this.children = new Set()
    this._task = task
    this._status = 'idle'
  }

  async onTrigger(context: C): Promise<void> {
    try {
      this._status = 'pending'
      await this._task.onTrigger(context)
      this._status = 'completed'
    } catch (error) {
      this._status = 'errored'
      await this._task.onError(context, error instanceof Error ? error : new Error(String(error)))
    }
  }

  async onReset(context: C): Promise<void> {
    if (this._status === 'idle') return
    try {
      this._status = 'pending'
      await this._task.onReset(context)
      this._status = 'idle'
    } catch (error) {
      this._status = 'errored'
      await this._task.onError(context, error instanceof Error ? error : new Error(String(error)))
    }
  }

  get status(): TaskStatus {
    return this._status
  }

  get id(): string {
    return this._task.id
  }
}

export class DagTaskScheduler<C = Record<string, unknown>> {
  private readonly _context: C
  private readonly _taskIdToTaskNode = new Map<string, TaskNode<C>>()
  private _built = false

  constructor(context: C) {
    this._context = context
  }

  /**
   * Add a task to the Dag.
   */
  add(task: ITask<C>): void {
    if (this._built) {
      this._report('Cannot add task after Dag is built')
      return
    }

    const { id } = task
    if (this._taskIdToTaskNode.has(id)) {
      this._report(`Task ${id} already exists`)
      return
    }

    const node = new TaskNode(task)
    this._taskIdToTaskNode.set(id, node)
  }

  /**
   * Build the Dag.
   *
   * @throws {Error} If a cycle is detected or a task depends on a non-existent task.
   */
  build(): void {
    if (this._built) {
      this._report('Dag is already built')
      return
    }

    this._buildGraph()
    this._topoSort()
    this._built = true
  }

  async trigger(id: string): Promise<void> {
    if (!this._built) {
      this._report('Dag is not built yet')
      return
    }

    const curNode = this._taskIdToTaskNode.get(id)
    if (!curNode) {
      this._report(`Task ${id} does not exist`)
      return
    }

    if (curNode.status === 'pending') {
      return
    }

    // A task cannot run until all of its dependencies have completed.
    if (!this._allDepsCompleted(curNode)) {
      return
    }

    // dont catch error here, let the caller handle it
    await curNode.onTrigger(this._context)
    await this._resetAllChildren(id)
    await this._tryTriggerNextTasks(id)
  }

  async run(): Promise<void> {
    const f = async (id: string) => {
      const node = this._taskIdToTaskNode.get(id)!
      await node.onTrigger(this._context)
    }

    const levels = this._topoSort()
    for (const level of levels) {
      // eslint-disable-next-line no-await-in-loop
      await Promise.all(level.map(f))
    }
  }

  private _buildGraph(): void {
    for (const [id, curNode] of this._taskIdToTaskNode.entries()) {
      for (const depId of curNode.deps) {
        if (!this._taskIdToTaskNode.has(depId)) {
          throw new Error(`Task ${id} depends on non-existent task ${depId}`)
        }

        const depNode = this._taskIdToTaskNode.get(depId)!
        depNode.children.add(id)
      }
    }
  }

  /**
   * Topological sort the graph.
   * @returns The topological order of the tasks organized by levels. Each level contains tasks that can be executed in parallel.
   * @throws {Error} If a cycle is detected in the graph.
   */
  private _topoSort(): string[][] {
    const inDegree = new Map<string, number>()
    for (const [id, node] of this._taskIdToTaskNode.entries()) {
      inDegree.set(id, node.deps.size)
    }

    const queue: TaskNode<C>[] = []
    for (const node of this._taskIdToTaskNode.values()) {
      if (node.deps.size === 0) {
        queue.push(node)
      }
    }

    let processedCount = 0
    const levels: string[][] = []

    while (queue.length > 0) {
      const preLevelSize = queue.length
      const curLevel: string[] = []

      for (let i = 0; i < preLevelSize; i++) {
        const cur = queue.shift()!
        curLevel.push(cur.id)
        processedCount++

        for (const child of cur.children.values()) {
          const childNode = this._taskIdToTaskNode.get(child)!
          inDegree.set(child, inDegree.get(child)! - 1)
          if (inDegree.get(child) === 0) {
            queue.push(childNode)
          }
        }
      }

      levels.push(curLevel)
    }

    const hasCycle = processedCount !== this._taskIdToTaskNode.size
    if (hasCycle) {
      throw new Error('Cycle detected')
    }

    return levels
  }

  private async _resetAllChildren(id: string): Promise<void> {
    const f = async (childId: string) => {
      const childNode = this._taskIdToTaskNode.get(childId)!
      await childNode.onReset(this._context)
      await this._resetAllChildren(childId)
    }

    const curNode = this._taskIdToTaskNode.get(id)!
    await Promise.all([...curNode.children].map(f)) // TODO: concurrency control
  }

  private async _tryTriggerNextTasks(id: string): Promise<void> {
    const f = async (childId: string) => {
      const childNode = this._taskIdToTaskNode.get(childId)!
      if (this._allDepsCompleted(childNode)) {
        await this.trigger(childId)
      }
    }

    const curNode = this._taskIdToTaskNode.get(id)!
    await Promise.all([...curNode.children].map(f)) // TODO: concurrency control
  }

  private _allDepsCompleted(node: TaskNode<C>): boolean {
    return [...node.deps].every(depId => {
      const dep = this._taskIdToTaskNode.get(depId)!
      return dep.status === 'completed'
    })
  }

  private _report(message: string): void {
    console.log(message)
  }
}
