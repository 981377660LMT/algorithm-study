/**
 * 将多个`布尔状态`压缩到一个整数中.
 *
 * @warning
 * 埋点上报等场景不要这样做，难以分析数据。需要将每个状态作为单独的字段上报，方便后续分析。
 */
export class StateCompressor {
  private state: number

  constructor(initialState = 0) {
    this.state = initialState
  }

  /**
   * 检查第 index 位是否为 true (1)
   * @param index 位索引 (从 0 开始)
   */
  has(index: number): boolean {
    return (this.state & (1 << index)) !== 0
  }

  /**
   * 将第 index 位设置为 true (1)
   * @param index 位索引
   */
  add(index: number): void {
    this.state |= 1 << index
  }

  /**
   * 将第 index 位设置为 false (0)
   * @param index 位索引
   */
  delete(index: number): void {
    this.state &= ~(1 << index)
  }

  /**
   * 切换第 index 位的状态 (0 -> 1, 1 -> 0)
   * @param index 位索引
   */
  toggle(index: number): void {
    this.state ^= 1 << index
  }

  /**
   * 检查是否包含所有指定的位
   */
  all(...indices: number[]): boolean {
    const mask = indices.reduce((pre, cur) => pre | (1 << cur), 0)
    return (this.state & mask) === mask
  }

  /**
   * 检查是否包含任意指定的位
   */
  any(...indices: number[]): boolean {
    const mask = indices.reduce((pre, cur) => pre | (1 << cur), 0)
    return (this.state & mask) !== 0
  }
}

{
  enum Permission {
    Read,
    Write,
    Execute
  }

  const permissions = new StateCompressor()

  permissions.add(Permission.Read)
  permissions.add(Permission.Write)

  console.log(permissions.has(Permission.Read)) // true
  console.log(permissions.has(Permission.Execute)) // false

  permissions.toggle(Permission.Write)
  console.log(permissions.has(Permission.Write)) // false

  permissions.add(Permission.Write)
  console.log(permissions.all(Permission.Read, Permission.Write)) // true
}
