```JS
// 获取父路径
const lastIndex = path.lastIndexOf('/')
return path.slice(0, lastIndex)

// 分割路径
private tranverse(filePath: string) {
    let root = this.root
    const pathList = filePath.split('/').filter(v => v)
    for (const name of pathList) {
      if (!root.children.has(name)) root.children.set(name, new File(''))
      root = root.children.get(name)!
    }
    return root
  }
```
