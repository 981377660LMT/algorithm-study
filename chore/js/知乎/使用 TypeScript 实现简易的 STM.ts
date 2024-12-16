/* eslint-disable no-lone-blocks */
type TransactionFunction = () => void

// 1. 定义事务和内存单元
class STM {
  private readonly memory: Map<string, any>
  private readonly transactions: TransactionFunction[]

  constructor() {
    this.memory = new Map()
    this.transactions = []
  }

  // 读取内存中的值
  read(key: string): any {
    return this.memory.get(key)
  }

  // 写入内存中的值
  write(key: string, value: any): void {
    this.memory.set(key, value)
  }

  // 添加一个事务
  addTransaction(transaction: TransactionFunction): void {
    this.transactions.push(transaction)
  }

  // 执行所有事务
  run(): void {
    for (const transaction of this.transactions) {
      transaction()
    }
    this.transactions.length = 0 // 清空事务队列
  }
}

// 2. 使用 STM 进行并发操作
{
  const stm = new STM()

  // 初始化共享变量
  stm.write('counter', 0)

  // 定义一个事务，增加计数器
  const incrementTransaction = () => {
    const current = stm.read('counter')
    stm.write('counter', current + 1)
  }

  // 添加多个事务
  for (let i = 0; i < 5; i++) {
    stm.addTransaction(incrementTransaction)
  }

  // 执行事务
  stm.run()

  // 读取最终结果
  console.log(stm.read('counter')) // 输出: 5
}

// c. 处理事务冲突（简化示例）
// 在实际的 STM 系统中，需要处理多个事务之间的冲突和重试。
class STMEnhanced extends STM {
  private readonly version: Map<string, number> // key -> version
  private readonly readSet: Map<string, number> // key -> version
  private readonly writeSet: Map<string, any> // key -> value

  constructor() {
    super()
    this.version = new Map()
    this.readSet = new Map()
    this.writeSet = new Map()
  }

  override read(key: string): any {
    const value = super.read(key)
    const ver = this.version.get(key) || 0
    this.readSet.set(key, ver)
    return value
  }

  override write(key: string, value: any): void {
    this.writeSet.set(key, value)
  }

  commit(): boolean {
    // 检查读集的版本是否未变
    for (const [key, ver] of this.readSet.entries()) {
      const currentVer = this.version.get(key) || 0
      if (currentVer !== ver) {
        return false // 版本变化，事务失败
      }
    }

    // 应用写集并更新版本
    for (const [key, value] of this.writeSet.entries()) {
      super.write(key, value)
      const currentVer = this.version.get(key) || 0
      this.version.set(key, currentVer + 1)
    }

    // 清空读集和写集
    this.readSet.clear()
    this.writeSet.clear()
    return true
  }

  runTransaction(transaction: TransactionFunction, maxRetries = 5): boolean {
    for (let i = 0; i < maxRetries; i++) {
      transaction()
      if (this.commit()) {
        return true
      }
      // 重试
      this.readSet.clear()
      this.writeSet.clear()
    }
    return false // 事务失败
  }
}

{
  const stmEnhanced = new STMEnhanced()
  stmEnhanced.write('counter', 0)

  // 定义一个可能导致冲突的事务
  const conflictingTransaction = () => {
    const current = stmEnhanced.read('counter')
    stmEnhanced.write('counter', current + 1)
  }

  // 并发执行多个事务
  for (let i = 0; i < 10; i++) {
    stmEnhanced.runTransaction(conflictingTransaction)
  }

  console.log(stmEnhanced.read('counter')) // 输出: 10
}
