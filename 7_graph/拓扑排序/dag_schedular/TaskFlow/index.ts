// 1. 按照层并发执行？
// 2.

export interface ITask<C> {
  deps: string[]
  onTrigger(context: C): void | Promise<void>
  onReset(context: C): void | Promise<void>
  onError(context: C, error: Error): void | Promise<void>
}

export type TaskStatus = 'idle' | 'running' | 'completed' | 'errored'

class TaskNode<C> {
  readonly id: string
  readonly task: ITask<C>
  readonly deps: Set<string>
  readonly children: Set<string>
  readonly state: TaskStatus = 'idle'

  constructor(id: string, task: ITask<C>) {
    this.id = id
    this.task = task
    this.deps = new Set(task.deps)
    this.children = new Set()
  }
}

export class TaskScheduler<C> {
  private readonly _context: C
  private readonly _nodes = new Map<string, TaskNode<C>>()

  constructor(context: C) {
    this._context = context
  }

  add(id: string, task: ITask<C>): void {
    if (this._nodes.has(id)) throw new Error(`Task ${id} already exists`)
    const node = new TaskNode(id, task)
    this._nodes.set(id, node)
  }

  build(): void {
    // // 更新依赖关系 & 检测循环
    // task.deps.forEach(depId => {
    //   if (!this._nodes.has(depId)) throw new Error(`Dependency ${depId} not found`)
    //   this._nodes.get(depId)!.children.add(id)
    // })
    // this.checkCycle(id) // 循环检测算法（如DFS）
    // !建图 + 环检测
  }

  async run(id: string) {
    // 初始触发无依赖的任务
    for (const node of this._nodes.values()) {
      if (node.deps.size === 0) {
        this.executeNode(node)
      }
    }
  }

  private async executeNode(node: TaskNode<C>) {
    if (node.state === TaskState.PENDING || node.state === TaskState.Finished) return

    // 检查所有依赖是否完成
    const depsCompleted = Array.from(node.deps).every(depId => this._nodes.get(depId)!.state === TaskState.Finished)

    if (!depsCompleted) return

    node.state = TaskState.PENDING
    try {
      await node.task.onTrigger(this._context)
      node.state = TaskState.Finished
      // 触发子节点检查
      node.children.forEach(childId => this.executeNode(this._nodes.get(childId)!))
    } catch (error) {
      node.state = TaskState.Errored
      node.task.onError(this._context, error)
    }
  }

  // 手动触发（示例）
  async triggerTask(id: string) {
    const node = this._nodes.get(id)
    if (!node) throw new Error('Task not found')

    // 重置当前任务及下游
    this.resetTask(node)
    await this.executeNode(node)
  }

  private resetTask(node: TaskNode<C>) {
    if (node.state === TaskState.RESET) return
    node.state = TaskState.IDLE
    node.task.onReset(this._context)
    node.children.forEach(childId => this.resetTask(this._nodes.get(childId)!))
  }
}
