type NamingStrategy = (baseName: string, suffix: number) => string

const defaultNamingStrategy: NamingStrategy = (baseName, suffix) => `${baseName}(${suffix})`

interface IUniqueNameGeneratorOptions {
  initialNames?: string[]
  namingStrategy?: NamingStrategy
}

/**
 * 一个支持增删改查的唯一名称生成器。
 * 删除名称后，该名称可以被后续操作复用。
 */
export class UniqueNameGenerator {
  private readonly _existingNames: Set<string> = new Set()
  private readonly _namingStrategy: NamingStrategy

  constructor(options: IUniqueNameGeneratorOptions = {}) {
    const { initialNames = [], namingStrategy = defaultNamingStrategy } = options
    this._namingStrategy = namingStrategy
    initialNames.forEach(name => this.add(name))
  }

  /**
   * 添加一个新名称。如果名称已存在，则从头开始查找最小的可用后缀来生成唯一版本。
   * @param name - 想要添加的基础名称。
   * @returns 分配的唯一名称。
   */
  add(name: string): string {
    if (!this._existingNames.has(name)) {
      this._existingNames.add(name)
      return name
    }

    let mex = 1
    let newName: string
    while (true) {
      newName = this._namingStrategy(name, mex)
      if (!this._existingNames.has(newName)) {
        break
      }
      mex++
    }

    this._existingNames.add(newName)
    return newName
  }

  has(name: string): boolean {
    return this._existingNames.has(name)
  }

  delete(name: string): boolean {
    return this._existingNames.delete(name)
  }

  /**
   * 重命名一个已存在的名称 (Update)。
   */
  rename(oldName: string, newName: string): string | undefined {
    if (!this.has(oldName)) {
      return undefined
    }
    this.delete(oldName)
    return this.add(newName)
  }

  getAllNames(): string[] {
    return Array.from(this._existingNames)
  }
}
