class File {
  isFile: boolean
  content: string
  children: Map<string, File>
  constructor(content: string) {
    this.isFile = false
    this.content = content
    this.children = new Map()
  }
}

class FileSystem {
  private root: File

  constructor() {
    this.root = new File('')
  }

  // 如果它是一个文件的路径，那么函数返回一个列表，仅包含这个文件的名字。
  // 如果它是一个文件夹的的路径，那么返回该 文件夹内 的所有文件和子文件夹的名字。
  // 你的返回结果（包括文件和子文件夹）应该按字典序排列。
  ls(path: string): string[] {
    let root = this.root
    const pathList = path.split('/').filter(v => v)
    for (const name of pathList) {
      root = root.children.get(name)! // 用户不会获取不存在文件的内容
    }
    // 注意'/'的情况
    if (root.isFile) return [pathList[pathList.length - 1]]
    else return [...root.children.keys()].sort()
  }

  mkdir(path: string): void {
    this.tranverse(path)
  }

  // 如果文件不存在，你需要创建包含给定文件内容的文件。
  // 如果文件已经存在，那么你需要将给定的文件内容 追加 在原本内容的后面
  addContentToFile(filePath: string, content: string): void {
    const root = this.tranverse(filePath)
    root.isFile = true
    root.content += content
  }

  // 输入 文件路径 ，以字符串形式返回该文件的 内容 。
  readContentFromFile(filePath: string): string {
    return this.tranverse(filePath).content
  }

  private tranverse(filePath: string) {
    let root = this.root
    const pathList = filePath.split('/').filter(v => v)
    for (const name of pathList) {
      if (!root.children.has(name)) root.children.set(name, new File(''))
      root = root.children.get(name)!
    }
    return root
  }

  static main() {
    const fs = new FileSystem()
    console.log(fs.ls('/')) // []
    fs.mkdir('/a/b/c')
    fs.addContentToFile('/a/b/c/d', 'hello')
    console.log(fs.ls('/')) // ['a']
    console.log(fs.readContentFromFile('/a/b/c/d')) // 'hello'
  }
}

FileSystem.main()

export {}
