// DAGScheduler

interface ITask<C> {
  readonly id: string
  readonly deps: string[]
  onTrigger(context: C): void | Promise<void>
  onReset(context: C): void | Promise<void>
  onError(context: C, error: Error): void | Promise<void>
}

type TaskStatus = 'idle' | 'pending' | 'completed' | 'errored'

class TaskNode<C> {
  readonly deps: Set<string>
  readonly children: Set<string>

  private readonly _task: ITask<C>

  private _status: TaskStatus
  private _locked: boolean

  constructor(task: ITask<C>) {
    this.deps = new Set(task.deps)
    this.children = new Set()
    this._task = task
    this._status = 'idle'
    this._locked = false
  }

  async onTrigger(context: C): Promise<void> {
    if (this._locked) return
    this._locked = true
    try {
      this._status = 'pending'
      await this._task.onTrigger(context)
      this._status = 'completed'
    } catch (error) {
      this._status = 'errored'
      await this._task.onError(context, error instanceof Error ? error : new Error(String(error)))
    } finally {
      this._locked = false
    }
  }

  async onReset(context: C): Promise<void> {
    if (this._locked) return
    if (this._status === 'idle') return
    this._locked = true
    try {
      this._status = 'pending'
      await this._task.onReset(context)
      this._status = 'idle'
    } catch (error) {
      this._status = 'errored'
      await this._task.onError(context, error instanceof Error ? error : new Error(String(error)))
    } finally {
      this._locked = false
    }
  }

  get status(): TaskStatus {
    return this._status
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
    this._verifyNoCyclesExist()
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
    await this._tryResetChildren(id)
    await this._tryTriggerNextTasks(id)
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

  private _verifyNoCyclesExist(): void {
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
    while (queue.length > 0) {
      const cur = queue.shift()!
      processedCount++
      for (const child of cur.children.values()) {
        const childNode = this._taskIdToTaskNode.get(child)!
        inDegree.set(child, inDegree.get(child)! - 1)
        if (inDegree.get(child) === 0) {
          queue.push(childNode)
        }
      }
    }

    const hasCycle = processedCount !== this._taskIdToTaskNode.size
    if (hasCycle) {
      throw new Error('Cycle detected')
    }
  }

  private async _tryResetChildren(id: string): Promise<void> {
    const f = async (childId: string) => {
      const childNode = this._taskIdToTaskNode.get(childId)!
      await childNode.onReset(this._context)
      await this._tryResetChildren(childId)
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
