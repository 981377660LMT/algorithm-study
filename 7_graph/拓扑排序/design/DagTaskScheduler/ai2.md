**设计思路** (By O1)

1. 用 Map 存储 task，包含依赖。
2. 利用拓扑排序（若有环则抛错）构建 DAG。
3. 提供 trigger(id) 调用，若依赖完成则执行任务；执行完成后尝试自动触发子节点。
4. reset(id) 递归重置当前节点及其子节点的状态。

```typescript
enum TaskState {
  Idle,
  Running,
  Completed,
  Error
}

export interface ITask<C> {
  deps: string[]
  onTrigger(context: C): void | Promise<void>
  onReset(context: C): void | Promise<void>
  onError(context: C, error: Error): void | Promise<void>
}

export class DAGTaskScheduler<C> {
  private tasks = new Map<string, { task: ITask<C>; state: TaskState }>()
  private graph = new Map<string, string[]>()
  private revGraph = new Map<string, string[]>()
  private built = false

  add(id: string, task: ITask<C>) {
    this.tasks.set(id, { task, state: TaskState.Idle })
    this.graph.set(id, task.deps || [])
    task.deps.forEach(dep => {
      if (!this.revGraph.has(dep)) this.revGraph.set(dep, [])
      this.revGraph.get(dep)!.push(id)
    })
  }

  build() {
    const visited = new Set<string>()
    const stack = new Set<string>()
    const dfs = (node: string) => {
      if (stack.has(node)) throw new Error(`DAG error: cycle at ${node}`)
      if (visited.has(node)) return
      stack.add(node)
      ;(this.graph.get(node) || []).forEach(dfs)
      stack.delete(node)
      visited.add(node)
    }
    for (const id of this.tasks.keys()) dfs(id)
    this.built = true
  }

  private canRun(id: string) {
    return (this.graph.get(id) || []).every(dep => {
      const depTask = this.tasks.get(dep)
      return depTask && depTask.state === TaskState.Completed
    })
  }

  async trigger(id: string, context: C) {
    if (!this.built) this.build()
    const item = this.tasks.get(id)
    if (!item) return
    if (item.state === TaskState.Error || item.state === TaskState.Running) return
    if (item.state === TaskState.Completed) {
      await this.reset(id, context)
    }
    if (this.canRun(id)) {
      item.state = TaskState.Running
      try {
        await item.task.onTrigger(context)
        item.state = TaskState.Completed
        const children = this.revGraph.get(id) || []
        for (const childId of children) {
          if (this.canRun(childId)) this.trigger(childId, context)
        }
      } catch (err) {
        item.state = TaskState.Error
        await item.task.onError(context, err as Error)
      }
    }
  }

  async reset(id: string, context: C) {
    const item = this.tasks.get(id)
    if (!item || item.state === TaskState.Idle) return
    item.state = TaskState.Idle
    await item.task.onReset(context)
    const children = this.revGraph.get(id) || []
    for (const childId of children) {
      await this.reset(childId, context)
    }
  }
}
```
