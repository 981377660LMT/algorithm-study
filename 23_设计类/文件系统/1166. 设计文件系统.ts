class FileSystem {
  private pathToValue: Map<string, number>
  constructor() {
    this.pathToValue = new Map()
    this.pathToValue.set('', -1)
  }

  // 重复路径fileSystemCreatePath需要报错
  createPath(path: string, value: number): boolean {
    if (this.pathToValue.has(path)) return false
    const parentPath = this.getParentPath(path)
    if (!this.pathToValue.has(parentPath)) return false
    this.pathToValue.set(path, value)
    return true
  }

  get(path: string): number {
    return this.pathToValue.get(path) || -1
  }

  private getParentPath(path: string): string {
    const lastIndex = path.lastIndexOf('/')
    return path.slice(0, lastIndex)
  }
}

const fs = new FileSystem()
fs.createPath('/leet', 1) // 返回 true
fs.createPath('/leet/code', 2) // 返回 true
console.log(fs.get('/leet/code')) // 返回 2
fs.createPath('/c/d', 1) // 返回 false 因为父路径 "/c" 不存在。
console.log(fs.get('/c')) // 返回 -1 因为该路径不存在。

export default 1
