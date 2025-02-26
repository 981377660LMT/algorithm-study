export interface ITask<C> {
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
    this._status = 'idle'
    this._task = task
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

export class DagTaskSchedular<C> {
  private readonly _context: C
  private readonly _taskIdToTaskNode = new Map<string, TaskNode<C>>()
  private _built = false

  constructor(context: C) {
    this._context = context
  }

  add(task: ITask<C>): void {
    if (this._built) {
      throw new Error('Cannot add task after Dag is built')
    }

    const { id } = task
    if (this._taskIdToTaskNode.has(id)) {
      throw new Error(`Task ${id} already exists`)
    }

    const node = new TaskNode(task)
    this._taskIdToTaskNode.set(id, node)
  }

  build(): void {
    if (this._built) {
      throw new Error('Dag is already built')
    }

    this._buildGraph()
    this._verifyNoCyclesExist()
    this._built = true
  }
  async run(id: string): Promise<void> {
    if (!this._built) {
      throw new Error('Dag is not built yet')
    }

    const curNode = this._taskIdToTaskNode.get(id)
    if (!curNode) {
      throw new Error(`Task ${id} does not exist`)
    }

    if (curNode.status === 'pending') {
      console.warn(`Cannot run task ${id}: task is pending`)
      return
    }

    // A task cannot run until all of its dependencies have completed.
    for (const depId of curNode.deps) {
      const depNode = this._taskIdToTaskNode.get(depId)!
      if (depNode.status !== 'completed') {
        console.log(`Cannot run task ${id}: dependency ${depId} is not completed`)
        return
      }
    }

    try {
      await curNode.onTrigger(this._context)
      await this._tryResetChildren(id)
      await this._tryTriggerNextTasks(id)
    } catch (error) {
      console.error(`Error in task ${id}`, error)
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
    await Promise.all([...curNode.children].map(f))
  }

  private async _tryTriggerNextTasks(id: string): Promise<void> {
    const f = async (childId: string) => {
      const childNode = this._taskIdToTaskNode.get(childId)!
      const allDepsCompleted = [...childNode.deps].every(depId => {
        const dep = this._taskIdToTaskNode.get(depId)!
        return dep.status === 'completed'
      })
      if (allDepsCompleted) {
        await this.run(childId)
      }
    }

    const curNode = this._taskIdToTaskNode.get(id)!
    await Promise.all([...curNode.children].map(f))
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const fetchUserData = async (): Promise<{ id: string }> => ({ id: '123' })
  const fetchPermissions = async (userId: string): Promise<string[]> => ['read', 'write']

  interface IContext {
    formData: Record<string, any>
  }
  const context: IContext = { formData: {} }
  const schedular = new DagTaskSchedular<IContext>(context)

  schedular.add({
    id: 'fetchUserData',
    deps: [],
    onTrigger: async ctx => {
      ctx.formData.user = await fetchUserData()
      console.log('user', ctx.formData.user)
    },
    onReset: ctx => {
      delete ctx.formData.user
      console.log('reset user')
    },
    onError: (ctx, err) => {
      console.error('Failed to fetch user data', err)
    }
  })

  schedular.add({
    id: 'fetchUserPermissions',
    deps: ['fetchUserData'],
    onTrigger: async ctx => {
      ctx.formData.permissions = await fetchPermissions(ctx.formData.user.id)
      console.log('permissions', ctx.formData.permissions)
    },
    onReset: ctx => {
      delete ctx.formData.permissions
      console.log('reset permissions')
    },
    onError: (ctx, err) => {
      console.error('Failed to fetch permissions', err)
    }
  })

  schedular.build()

  schedular
    .run('fetchUserData')
    .then(() => {
      console.log('done')
    })
    .catch(err => {
      console.error(err)
    })
}
