export {}

/**
 * 策略接口：定义如何从一个项目中提取其唯一标识符（Key）。
 * @template T 项目的类型。
 * @template K 唯一标识符的类型。
 */
interface IIdentityStrategy<T, K> {
  getKey(item: T): K
}

/**
 * 策略接口：定义当冲突发生时，如何生成一系列候选项目。
 * 使用生成器函数可以优雅地处理有状态的序列生成。
 * @template T 项目的类型。
 */
interface IResolutionStrategy<T> {
  resolve(item: T): IterableIterator<T>
}

/**
 * 存储接口：定义了唯一项目的存储和检索行为。
 * 这将核心逻辑与持久化机制（内存、数据库等）解耦。
 * @template T 项目的类型。
 * @template K 唯一标识符的类型。
 */
interface IUniqueItemStore<T, K> {
  has(key: K): boolean
  add(key: K, item: T): void
  delete(key: K): boolean
  getAll(): T[]
}

/**
 * 一个基于内存 Map 的 IUniqueItemStore 的默认实现。
 */
class InMemoryStore<T, K> implements IUniqueItemStore<T, K> {
  private store = new Map<K, T>()
  has(key: K): boolean {
    return this.store.has(key)
  }
  add(key: K, item: T): void {
    this.store.set(key, item)
  }
  delete(key: K): boolean {
    return this.store.delete(key)
  }
  getAll(): T[] {
    return Array.from(this.store.values())
  }
}

/**
 * 一个高度抽象的、可组合的唯一项生成器。
 * 它协调标识、解析和存储策略来完成工作。
 */
class Generator<T, K> {
  constructor(
    private identityStrategy: IIdentityStrategy<T, K>,
    private resolutionStrategy: IResolutionStrategy<T>,
    private store: IUniqueItemStore<T, K>
  ) {}

  /**
   * 添加一个项目，如果冲突则通过解析策略生成唯一的版本。
   * @param item 要添加的原始项目。
   * @returns 最终被存储的唯一项目。
   */
  add(item: T): T {
    const initialKey = this.identityStrategy.getKey(item)

    if (!this.store.has(initialKey)) {
      this.store.add(initialKey, item)
      return item
    }

    // 冲突发生，启动解析策略的生成器
    for (const candidate of this.resolutionStrategy.resolve(item)) {
      const candidateKey = this.identityStrategy.getKey(candidate)
      if (!this.store.has(candidateKey)) {
        this.store.add(candidateKey, candidate)
        return candidate
      }
    }

    // 如果策略的生成器耗尽了所有候选项仍然冲突，则抛出错误。
    throw new Error('Resolution strategy failed to produce a unique item.')
  }

  has(item: T): boolean {
    return this.store.has(this.identityStrategy.getKey(item))
  }

  delete(item: T): boolean {
    return this.store.delete(this.identityStrategy.getKey(item))
  }

  getAll(): T[] {
    return this.store.getAll()
  }
}

// --- 使用示例 ---

// 示例：重新实现文件名生成器

// 1. 标识策略：文件名本身就是它的 Key
class FileNameIdentity implements IIdentityStrategy<string, string> {
  getKey(item: string): string {
    return item
  }
}

// 2. 解析策略：使用生成器函数来创建 "name(k)" 序列
class SuffixResolver implements IResolutionStrategy<string> {
  // 使用 Map 来为每个基础名称维护下一个后缀索引
  private nextSuffix = new Map<string, number>();

  *resolve(item: string): IterableIterator<string> {
    let k = this.nextSuffix.get(item) || 1
    while (true) {
      const candidate = `${item}(${k})`
      this.nextSuffix.set(item, k + 1) // 预先增加 k，以便下次从正确的位置开始
      k++
      yield candidate
    }
  }
}

// 3. 组合使用
const nameGenerator = new Generator(
  new FileNameIdentity(),
  new SuffixResolver(),
  new InMemoryStore<string, string>()
)

;['doc', 'image', 'doc', 'doc'].forEach(name => nameGenerator.add(name))
console.log(nameGenerator.getAll()) // [ 'doc', 'image', 'doc(1)', 'doc(2)' ]
console.log(nameGenerator.add('image')) // 'image(1)'
