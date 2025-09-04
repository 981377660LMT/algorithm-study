// TODO: state、view

/**
 * 变更集接口，用于记录在一次操作中发生的所有变更意图。
 */
interface ChangeSet {
  // 记录需要更新内容的节点 (nodeId -> newContent)
  contentUpdates: Map<string, string>
  // 记录其他一次性的写操作
  genericWrites: (() => void)[]
  // 记录需要读取布局信息并执行回调的请求
  layoutReads: { nodeId: string; callback: (rect: DOMRect) => void }[]
}

/**
 * 实现事务性批量更新算法的 UI 更新器。
 */
class UIBatchUpdater {
  private opDepth = 0 // 使用深度计数来处理嵌套，比布尔值更健壮
  private activeChangeSet: ChangeSet | null = null

  /**
   * 算法的核心入口：执行一个操作单元。
   * @param workload 包含一系列修改逻辑的函数
   */
  operation<T>(workload: () => T): T {
    // 嵌套守卫：如果已在操作中，直接执行即可
    if (this.isInOperation()) {
      return workload()
    }

    this.setUpOperation()
    try {
      return workload()
    } finally {
      this.finishOperation()
    }
  }

  // --- 模拟的修改方法 ---

  /**
   * 记录一个更新节点内容的意图。
   * @param nodeId 目标节点的 ID
   * @param content 新的内容
   */
  updateContent(nodeId: string, content: string) {
    if (!this.isInOperation()) {
      throw new Error('Cannot perform updates outside of an operation.')
    }
    console.log(`    (Queued: update content for #${nodeId})`)
    this.activeChangeSet!.contentUpdates.set(nodeId, content)
  }

  /**
   * 记录一个读取节点布局并执行回调的意图。
   * @param nodeId 目标节点的 ID
   * @param callback 获取到布局信息后要执行的回调
   */
  readLayout(nodeId: string, callback: (rect: DOMRect) => void) {
    if (!this.isInOperation()) {
      throw new Error('Cannot perform reads outside of an operation.')
    }
    console.log(`    (Queued: read layout for #${nodeId})`)
    this.activeChangeSet!.layoutReads.push({ nodeId, callback })
  }

  /**
   * 检查当前是否在一个活动的操作中。
   */
  private isInOperation(): boolean {
    return this.opDepth > 0
  }

  /**
   * 准备开始一个新操作（如果尚未开始）。
   */
  private setUpOperation() {
    this.opDepth++
    if (this.opDepth === 1) {
      console.log('>>> Transaction Started')
      this.activeChangeSet = {
        contentUpdates: new Map(),
        layoutReads: [],
        genericWrites: []
      }
    }
  }

  /**
   * 结束操作并提交所有变更。
   */
  private finishOperation() {
    this.opDepth--
    // 只有最外层的操作结束时才真正提交
    if (this.opDepth === 0) {
      console.log('<<< Transaction Finishing...')
      this.commit(this.activeChangeSet!)
      this.activeChangeSet = null
      console.log('<<< Transaction Committed & Finished')
    }
  }

  /**
   * 提交变更集，按优化的“写 -> 读 -> 回调”顺序执行。
   * @param changeSet 要提交的变更集
   */
  private commit(changeSet: ChangeSet) {
    // --- 1. 批量写入阶段 (Batch Write Phase) ---
    // 在此阶段执行所有不依赖于布局信息的 DOM 修改。
    console.log('  |- 1. WRITE PHASE: Applying content and generic updates...')
    changeSet.contentUpdates.forEach((content, nodeId) => {
      const el = document.getElementById(nodeId)
      if (el) {
        console.log(`     - Updating node #${nodeId} content to "${content}"`)
        el.textContent = content
      }
    })
    changeSet.genericWrites.forEach(writeAction => writeAction())

    // --- 2. 批量读取阶段 (Batch Read Phase) ---
    // 在此阶段一次性地执行所有 DOM 布局读取，以避免强制同步布局。
    console.log('  |- 2. READ PHASE: Reading layout information...')
    const readResults = changeSet.layoutReads.map(({ nodeId }) => {
      const el = document.getElementById(nodeId)
      if (el) {
        console.log(`     - Reading layout of node #${nodeId}`)
        // 强制同步布局在这里发生一次（如果需要）
        return el.getBoundingClientRect()
      }
      return null
    })

    // --- 3. 回调/收尾写入阶段 (Callback/Finalize Write Phase) ---
    // 在此阶段执行依赖于读取结果的操作。
    console.log('  |- 3. CALLBACK PHASE: Using read results...')
    changeSet.layoutReads.forEach(({ callback }, index) => {
      const rect = readResults[index]
      if (rect) {
        callback(rect)
      }
    })
  }
}

export {}

if (require.main === module) {
  const updater = new UIBatchUpdater()

  function handleUserTyping() {
    // 这是一个顶层操作，模拟用户输入
    updater.operation(() => {
      console.log("  Workload: User typed 'Hello'. Updating text1.")
      updater.updateContent('text1', 'Hello')

      // 模拟一个由输入触发的自动操作，这是一个嵌套操作
      autoCompleteBrackets()

      console.log('  Workload: User moved cursor. Reading text1 layout to position cursor.')
      // 读取布局以更新光标位置
      updater.readLayout('text1', rect => {
        const cursorEl = document.getElementById('cursor')
        if (cursorEl) {
          console.log(
            `     - Positioning cursor based on text1 layout (top: ${rect.top}, left: ${rect.left})`
          )
          cursorEl.style.top = `${rect.top}px`
          cursorEl.style.left = `${rect.left + rect.width}px`
        }
      })
    })
  }

  function autoCompleteBrackets() {
    // 这个函数被嵌套调用，它也启动一个 operation
    // 但由于嵌套守卫，它不会立即提交，而是将变更合并到外部操作中
    updater.operation(() => {
      console.log('  Nested Workload: Auto-completing brackets in text2.')
      updater.updateContent('text2', 'Content with ()')
    })
  }

  // 运行示例
  handleUserTyping()
}
