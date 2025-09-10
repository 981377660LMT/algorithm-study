/**
 * 定义了当名称冲突时生成新名称的策略函数。
 * @param baseName - 原始名称。
 * @param suffix - 当前建议的后缀数字。
 * @returns 一个新的、建议的唯一名称。
 */
type NamingStrategy = (baseName: string, suffix: number) => string

/**
 * 默认的命名策略，生成如 "name(1)", "name(2)" 的格式。
 */
const defaultNamingStrategy: NamingStrategy = (baseName, suffix) => `${baseName}(${suffix})`

interface IUniqueNameGeneratorAddOnlyOptions {
  initialNames?: string[]
  namingStrategy?: NamingStrategy
}

/**
 * 一个高效的唯一名称生成类，允许自定义重名规则。
 */
export class UniqueNameGeneratorAddOnly {
  /** 存储每个基础名称下一个可用的后缀数字。避免了在冲突时从 1 开始的重复搜索。 */
  private readonly _nextSuffix: Map<string, number> = new Map()
  private readonly _namingStrategy: NamingStrategy

  constructor(options: IUniqueNameGeneratorAddOnlyOptions = {}) {
    const { initialNames = [], namingStrategy = defaultNamingStrategy } = options
    this._namingStrategy = namingStrategy
    initialNames.forEach(name => this.add(name))
  }

  add(name: string): string {
    if (!this._nextSuffix.has(name)) {
      this._nextSuffix.set(name, 1)
      return name
    }
    let suffix = this._nextSuffix.get(name)!
    let newName = this._namingStrategy(name, suffix)
    while (this._nextSuffix.has(newName)) {
      suffix++
      newName = this._namingStrategy(name, suffix)
    }
    this._nextSuffix.set(name, suffix + 1)
    this._nextSuffix.set(newName, 1)
    return newName
  }
}

// https://leetcode.cn/problems/making-file-names-unique/description/
function getFolderNames(names: string[]): string[] {
  const generator = new UniqueNameGeneratorAddOnly({
    namingStrategy: (baseName, suffix) => `${baseName}(${suffix})`
  })
  return names.map(name => generator.add(name))
}

// --- 使用示例 ---
// 1. 使用默认策略 "name(k)"
const defaultGenerator = new UniqueNameGeneratorAddOnly({ initialNames: ['doc', 'image', 'doc'] })
console.log(defaultGenerator.add('doc')) // 'doc(2)'

// 2. 使用自定义下划线策略 "name_k"
const underscoreStrategy: NamingStrategy = (base, k) => `${base}_${k}`
const customGenerator = new UniqueNameGeneratorAddOnly({
  initialNames: ['file', 'data', 'file'],
  namingStrategy: underscoreStrategy
})
console.log(customGenerator.add('file')) // 'file_2'
