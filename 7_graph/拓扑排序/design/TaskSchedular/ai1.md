# DAG任务调度器设计与实现 (By Claude 3.7 Sonnect Thinking)

## 需求分析

这个DAG（有向无环图）任务调度器主要用于前端表单数据源联动场景，需要支持：

1. 构建任务依赖关系图
2. 手动/自动触发任务执行
3. 依赖任务变化时重置子任务
4. 错误处理机制

## 设计思路

1. **任务状态管理**：跟踪每个任务的状态（待执行、执行中、已完成、错误、已重置）
2. **DAG构建**：使用邻接表表示任务依赖关系
3. **执行机制**：检查依赖任务状态，按拓扑顺序执行
4. **重置机制**：当任务重新执行时，级联重置所有依赖它的任务

## 代码实现

```typescript
// 任务状态枚举
enum TaskState {
  PENDING = 'PENDING',
  RUNNING = 'RUNNING',
  COMPLETED = 'COMPLETED',
  ERROR = 'ERROR',
  RESET = 'RESET'
}

// 内部任务接口，扩展公共接口
interface InternalTask<C> extends ITask<C> {
  id: string
  state: TaskState
  children: string[] // 依赖于此任务的子任务
}

export class DAGTaskScheduler<C> {
  private tasks: Map<string, InternalTask<C>> = new Map()
  private context: C
  private built: boolean = false

  constructor(context: C) {
    this.context = context
  }

  /**
   * 添加任务到调度器
   */
  add(id: string, task: ITask<C>): DAGTaskScheduler<C> {
    if (this.built) {
      throw new Error('Cannot add task after DAG is built')
    }

    if (this.tasks.has(id)) {
      throw new Error(`Task with id ${id} already exists`)
    }

    const internalTask: InternalTask<C> = {
      ...task,
      id,
      state: TaskState.PENDING,
      children: []
    }

    this.tasks.set(id, internalTask)
    return this
  }

  /**
   * 构建DAG，设置依赖关系并检测环
   */
  build(): DAGTaskScheduler<C> {
    if (this.built) {
      return this
    }

    // 1. 建立依赖关系
    for (const [id, task] of this.tasks.entries()) {
      for (const depId of task.deps) {
        if (!this.tasks.has(depId)) {
          throw new Error(`Task ${id} depends on non-existent task ${depId}`)
        }
        // 将当前任务添加为依赖任务的子任务
        const depTask = this.tasks.get(depId)!
        depTask.children.push(id)
      }
    }

    // 2. 检测环
    this.detectCycle()

    this.built = true
    return this
  }

  /**
   * 检测DAG中是否存在环
   */
  private detectCycle(): void {
    const visited: Map<string, number> = new Map()
    // 0: 未访问, 1: 访问中, 2: 已完成访问

    const dfs = (id: string): boolean => {
      const visitStatus = visited.get(id)
      if (visitStatus === 1) return true // 发现环
      if (visitStatus === 2) return false // 已处理过的节点

      visited.set(id, 1) // 标记为访问中

      const task = this.tasks.get(id)!
      for (const childId of task.children) {
        if (dfs(childId)) {
          return true // 发现环
        }
      }

      visited.set(id, 2) // 标记为已完成
      return false
    }

    // 检查所有节点
    for (const id of this.tasks.keys()) {
      if (!visited.has(id) && dfs(id)) {
        throw new Error(`Cycle detected in task dependencies involving task ${id}`)
      }
    }
  }

  /**
   * 运行指定任务
   */
  async run(id: string): Promise<void> {
    if (!this.built) {
      throw new Error('DAG is not built yet')
    }

    const task = this.tasks.get(id)
    if (!task) {
      throw new Error(`Task ${id} does not exist`)
    }

    // 检查所有依赖任务是否已完成
    for (const depId of task.deps) {
      const depTask = this.tasks.get(depId)!
      if (depTask.state !== TaskState.COMPLETED) {
        console.log(`Cannot run task ${id}: dependency ${depId} is not completed`)
        return
      }
    }

    // 重置所有子任务
    await this.resetChildren(id)

    // 执行任务
    try {
      task.state = TaskState.RUNNING
      await Promise.resolve(task.onTrigger(this.context))
      task.state = TaskState.COMPLETED

      // 触发可能的后续任务
      await this.triggerNextTasks(id)
    } catch (error) {
      task.state = TaskState.ERROR
      await Promise.resolve(
        task.onError(this.context, error instanceof Error ? error : new Error(String(error)))
      )
    }
  }

  /**
   * 重置任务及其所有子任务
   */
  private async resetChildren(id: string): Promise<void> {
    const task = this.tasks.get(id)!

    for (const childId of task.children) {
      const childTask = this.tasks.get(childId)!

      if (childTask.state !== TaskState.PENDING) {
        childTask.state = TaskState.RESET

        try {
          await Promise.resolve(childTask.onReset(this.context))
          childTask.state = TaskState.PENDING
        } catch (error) {
          childTask.state = TaskState.ERROR
          await Promise.resolve(
            childTask.onError(
              this.context,
              error instanceof Error ? error : new Error(String(error))
            )
          )
        }

        // 递归重置子任务的所有子任务
        await this.resetChildren(childId)
      }
    }
  }

  /**
   * 尝试触发后续任务
   */
  private async triggerNextTasks(id: string): Promise<void> {
    const task = this.tasks.get(id)!

    for (const childId of task.children) {
      const childTask = this.tasks.get(childId)!

      // 检查子任务的所有依赖是否已完成
      const allDepsCompleted = childTask.deps.every(
        depId => this.tasks.get(depId)!.state === TaskState.COMPLETED
      )

      // 如果所有依赖都已完成且当前任务不是运行中，自动执行该任务
      if (allDepsCompleted && childTask.state !== TaskState.RUNNING) {
        await this.run(childId)
      }
    }
  }

  /**
   * 重置指定任务
   */
  async reset(id: string): Promise<void> {
    if (!this.built) {
      throw new Error('DAG is not built yet')
    }

    const task = this.tasks.get(id)
    if (!task) {
      throw new Error(`Task ${id} does not exist`)
    }

    try {
      task.state = TaskState.RESET
      await Promise.resolve(task.onReset(this.context))
      task.state = TaskState.PENDING

      // 重置所有子任务
      await this.resetChildren(id)
    } catch (error) {
      task.state = TaskState.ERROR
      await Promise.resolve(
        task.onError(this.context, error instanceof Error ? error : new Error(String(error)))
      )
    }
  }

  /**
   * 获取任务状态
   */
  getTaskState(id: string): TaskState | undefined {
    const task = this.tasks.get(id)
    return task?.state
  }

  /**
   * 更新上下文
   */
  setContext(context: C): void {
    this.context = context
  }
}
```

## 使用示例

```typescript
// 创建调度器实例
const context = { formData: {} }
const scheduler = new DAGTaskScheduler(context)

// 添加任务
scheduler.add('fetchUserData', {
  deps: [],
  onTrigger: async ctx => {
    ctx.formData.user = await fetchUserData()
  },
  onReset: ctx => {
    delete ctx.formData.user
  },
  onError: (ctx, err) => {
    console.error('Failed to fetch user data', err)
  }
})

scheduler.add('fetchUserPermissions', {
  deps: ['fetchUserData'],
  onTrigger: async ctx => {
    ctx.formData.permissions = await fetchPermissions(ctx.formData.user.id)
  },
  onReset: ctx => {
    delete ctx.formData.permissions
  },
  onError: (ctx, err) => {
    console.error('Failed to fetch permissions', err)
  }
})

// 构建DAG
scheduler.build()

// 执行任务
await scheduler.run('fetchUserData')
```

这个实现支持了任务依赖关系管理、状态跟踪、自动触发和错误处理，适合处理前端表单数据源联动场景。

---

# DAGTaskScheduler 优化方案与关键问题

## 优化方案

### 1. 性能优化

- **并行任务执行**：允许没有相互依赖的任务并行执行
- **任务批处理**：相似任务分组执行，减少上下文切换
- **懒加载机制**：任务定义与执行分离，需要时才实例化执行环境
- **增量构建DAG**：允许动态添加任务而不需要完全重建

### 2. 功能增强

- **任务优先级**：支持指定任务优先级以优化执行顺序
- **循环检测优化**：实现更高效的环检测算法
- **部分执行**：支持从DAG中选择性执行部分子图
- **状态持久化**：支持保存/恢复执行状态

### 3. 开发体验提升

- **可视化调试**：提供DAG可视化展示工具
- **任务监控面板**：实时监控任务状态和执行时间
- **事件系统**：支持任务生命周期事件订阅

## 关键问题与回答

### 问题1：如何处理大规模DAG中的性能问题？

**回答**：

- 实现任务分块处理，将大DAG拆分为多个子图
- 使用Web Workers实现并行执行无依赖任务
- 引入惰性执行机制，仅在需要结果时才执行任务
- 缓存中间结果，避免重复计算

### 问题2：如何处理长时间运行任务？

**回答**：

- 实现任务超时机制，可配置最大执行时间
- 支持任务取消操作，中断执行中的任务
- 提供进度报告接口，长任务周期性报告完成进度
- 支持任务断点续执行

### 问题3：如何增强错误处理和恢复机制？

**回答**：

- 实现细粒度错误分类，区分致命错误和非致命错误
- 支持自定义恢复策略，包括重试、跳过或回滚
- 提供错误传播控制，决定错误是否影响下游任务
- 维护执行快照，支持从故障点恢复

### 问题4：如何优化表单场景下的联动体验？

**回答**：

- 引入节流和防抖机制，避免频繁触发连锁反应 => `这个也可以由业务方自己控制，框架不关注，因为每个任务的触发是由业务方自己控制的，框架内置可能适得其反`
- 实现预计算，根据历史数据预测和准备可能的数据
- 支持增量更新，只处理变化的数据部分
- 优化渲染周期，将多个状态变更合并更新UI

### 问题5：如何使系统更易于维护和扩展？

**回答**：

- 采用插件架构，支持功能模块化扩展
- 提供标准化的任务生命周期钩子
- 实现中间件系统，支持横切关注点处理
- 完善测试工具，支持任务模拟和场景重现

这些优化可以根据具体项目需求和资源约束进行选择性实施，逐步提升系统性能和可用性。
