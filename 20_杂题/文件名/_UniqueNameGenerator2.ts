export {}

interface IEqualityComparer<T> {
  equals(a: T, b: T): boolean
  getHashCode(item: T): string | number
}

class DefaultEqualityComparer<T> implements IEqualityComparer<T> {
  equals(a: T, b: T): boolean {
    return a === b
  }

  getHashCode(item: T): string | number {
    // 对于对象，这将是 "[object Object]"，因此对于复杂对象需要自定义比较器。
    return String(item)
  }
}

/**
 * 定义了当项目冲突时如何生成新的候选项目。
 * 这是一个有状态的接口，允许解析器自身维护状态。
 */
interface IConflictResolver<T> {
  /**
   * @param item - 发生冲突的原始项目。
   * @returns 一个新的、建议的候选项目。
   */
  resolve(item: T): T
  /**
   * 当一个项目被成功添加后，可以调用此方法来重置或更新解析器的内部状态。
   * @param baseItem - 冲突的原始项目。
   * @param resolvedItem - 最终被唯一化并添加的项目。
   */
  onResolved(baseItem: T, resolvedItem: T): void
}

/**
 * 一个高度通用和抽象的唯一项生成和管理类。
 * @template T 要管理的项目的类型。
 */
class UniqueGenerator<T> {
  private _existingItems: Map<string | number, T>
  private _comparer: IEqualityComparer<T>
  private _resolver: IConflictResolver<T>

  constructor(
    resolver: IConflictResolver<T>,
    comparer: IEqualityComparer<T> = new DefaultEqualityComparer<T>(),
    initialItems: T[] = []
  ) {
    this._resolver = resolver
    this._comparer = comparer
    this._existingItems = new Map()

    initialItems.forEach(item => this.add(item))
  }

  /**
   * 添加一个新项目。如果项目已存在，则使用冲突解决策略生成一个唯一的版本。
   * @param item - 想要添加的项目。
   * @returns 最终被添加到集合中的唯一项目。
   */
  add(item: T): T {
    let current = item
    let hashCode = this._comparer.getHashCode(current)

    while (this._existingItems.has(hashCode)) {
      const existingItem = this._existingItems.get(hashCode)!
      if (this._comparer.equals(current, existingItem)) {
        // 发生冲突，请求解析器生成新候选项
        current = this._resolver.resolve(item)
        hashCode = this._comparer.getHashCode(current)
        continue
      }
      // 哈希冲突，但项目不相等（在自定义比较器中可能发生），继续检查
      break
    }

    this._existingItems.set(hashCode, current)
    this._resolver.onResolved(item, current)
    return current
  }

  has(item: T): boolean {
    const hashCode = this._comparer.getHashCode(item)
    return this._existingItems.has(hashCode)
  }

  delete(item: T): boolean {
    const hashCode = this._comparer.getHashCode(item)
    return this._existingItems.delete(hashCode)
  }

  getAll(): T[] {
    return Array.from(this._existingItems.values())
  }
}

// --- 使用示例 ---

// 示例1：实现一个用于字符串的经典 "(k)" 后缀解析器
class FileNameResolver implements IConflictResolver<string> {
  private _nextSuffix = new Map<string, number>()

  resolve(item: string): string {
    const k = this._nextSuffix.get(item) || 1
    return `${item}(${k})`
  }

  onResolved(baseItem: string, resolvedItem: string): void {
    const k = this._nextSuffix.get(baseItem) || 1
    if (baseItem !== resolvedItem) {
      this._nextSuffix.set(baseItem, k + 1)
    }
    this._nextSuffix.set(resolvedItem, 1)
  }
}

const nameGenerator = new UniqueGenerator(new FileNameResolver(), undefined, [
  'doc',
  'image',
  'doc'
])
console.log(nameGenerator.getAll()) // [ 'doc', 'image', 'doc(1)' ]
console.log(nameGenerator.add('doc')) // 'doc(2)'

// 示例2：管理具有唯一 `id` 的对象
interface User {
  id: number
  name: string
}

class UserIdComparer implements IEqualityComparer<User> {
  equals(a: User, b: User): boolean {
    return a.id === b.id
  }
  getHashCode(item: User): number {
    return item.id
  }
}

class UserConflictResolver implements IConflictResolver<User> {
  private _nextId = 100 // 假设从100开始分配新ID

  resolve(item: User): User {
    // 冲突时，忽略原始ID，分配一个全新的ID
    return { ...item, id: this._nextId }
  }

  onResolved(baseItem: User, resolvedItem: User): void {
    if (baseItem.id !== resolvedItem.id) {
      this._nextId++
    }
  }
}

const user1: User = { id: 1, name: 'Alice' }
const user2: User = { id: 2, name: 'Bob' }
const user3: User = { id: 1, name: 'Alice Clone' } // ID冲突

const userGenerator = new UniqueGenerator(new UserConflictResolver(), new UserIdComparer(), [
  user1,
  user2,
  user3
])
console.log(userGenerator.getAll())
// [
//   { id: 1, name: 'Alice' },
//   { id: 2, name: 'Bob' },
//   { id: 100, name: 'Alice Clone' }
// ]
