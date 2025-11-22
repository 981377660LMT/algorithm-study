这是一个非常硬核的系统设计题目。**虚拟文件系统 (VFS)** 是操作系统（如 Linux）和现代编辑器（如 VS Code）的核心基石。

它的核心目标是：**提供统一的 API（open, read, write），屏蔽底层物理存储的差异（内存、磁盘、网络、Git、Zip）。**

我们将实现一个功能完备的 VFS，包含以下核心特性：

1.  **统一接口**：`IFileSystem` 抽象基类。
2.  **挂载机制 (Mounting)**：将不同的文件系统挂载到路径树上（例如把 `MemoryFS` 挂载到 tmp，把 `NetworkFS` 挂载到 `/remote`）。
3.  **路径解析**：自动将 `/tmp/a.txt` 路由到对应的文件系统实例。
4.  **能力检测**：区分只读/读写文件系统。

---

### 1. 核心类型定义

首先定义文件系统的标准接口。这决定了所有具体实现必须遵守的契约。

```typescript
// 文件类型
export enum FileType {
  Unknown = 0,
  File = 1,
  Directory = 2,
  Symlink = 64
}

// 文件状态信息
export interface FileStat {
  type: FileType
  ctime: number
  mtime: number
  size: number
}

// 核心抽象接口：任何文件系统（内存、磁盘、网络）都要实现它
export interface IFileSystem {
  // 基础能力
  readonly: boolean

  // 核心操作
  stat(path: string): Promise<FileStat>
  readFile(path: string): Promise<Uint8Array>
  writeFile(path: string, content: Uint8Array): Promise<void>
  readDirectory(path: string): Promise<[string, FileType][]>
  createDirectory(path: string): Promise<void>
  delete(path: string, recursive?: boolean): Promise<void>
  rename(oldPath: string, newPath: string): Promise<void>
}

// 错误定义
export class FileSystemError extends Error {
  constructor(public code: string, message: string) {
    super(`[${code}] ${message}`)
  }
  static FileNotFound(path: string) {
    return new FileSystemError('ENOENT', `No such file or directory: ${path}`)
  }
  static FileExists(path: string) {
    return new FileSystemError('EEXIST', `File exists: ${path}`)
  }
  static IsDirectory(path: string) {
    return new FileSystemError('EISDIR', `Is a directory: ${path}`)
  }
  static NotDirectory(path: string) {
    return new FileSystemError('ENOTDIR', `Not a directory: ${path}`)
  }
  static NoPermissions(path: string) {
    return new FileSystemError('EACCES', `Permission denied: ${path}`)
  }
}
```

### 2. 实现一个具体的后端：内存文件系统 (MemoryFS)

这是最常用的 VFS 实现，用于测试或临时存储。我们需要用一个树状结构来模拟目录。

```typescript
// 简单的内存节点结构
interface MemoryNode {
  type: FileType
  name: string
  content?: Uint8Array // 只有文件有内容
  children?: Map<string, MemoryNode> // 只有目录有子节点
  mtime: number
}

export class MemoryFileSystem implements IFileSystem {
  readonly = false
  private root: MemoryNode

  constructor() {
    this.root = { type: FileType.Directory, name: '', children: new Map(), mtime: Date.now() }
  }

  // --- 辅助方法：解析路径找到节点 ---
  private findNode(path: string): MemoryNode | undefined {
    if (path === '/' || path === '') return this.root
    const parts = path.split('/').filter(p => p.length > 0)
    let current = this.root

    for (const part of parts) {
      if (current.type !== FileType.Directory || !current.children) {
        return undefined
      }
      const next = current.children.get(part)
      if (!next) return undefined
      current = next
    }
    return current
  }

  // --- 接口实现 ---

  async stat(path: string): Promise<FileStat> {
    const node = this.findNode(path)
    if (!node) throw FileSystemError.FileNotFound(path)
    return {
      type: node.type,
      ctime: node.mtime,
      mtime: node.mtime,
      size: node.content ? node.content.byteLength : 0
    }
  }

  async readFile(path: string): Promise<Uint8Array> {
    const node = this.findNode(path)
    if (!node) throw FileSystemError.FileNotFound(path)
    if (node.type === FileType.Directory) throw FileSystemError.IsDirectory(path)
    return node.content || new Uint8Array(0)
  }

  async writeFile(path: string, content: Uint8Array): Promise<void> {
    const parts = path.split('/').filter(p => p.length > 0)
    const fileName = parts.pop()!
    const dirPath = '/' + parts.join('/')

    const dirNode = this.findNode(dirPath)
    if (!dirNode) throw FileSystemError.FileNotFound(dirPath)
    if (dirNode.type !== FileType.Directory) throw FileSystemError.NotDirectory(dirPath)

    // 写入或覆盖
    const now = Date.now()
    dirNode.children!.set(fileName, {
      type: FileType.File,
      name: fileName,
      content: content,
      mtime: now
    })
    dirNode.mtime = now // 更新父目录时间
  }

  async readDirectory(path: string): Promise<[string, FileType][]> {
    const node = this.findNode(path)
    if (!node) throw FileSystemError.FileNotFound(path)
    if (node.type !== FileType.Directory) throw FileSystemError.NotDirectory(path)

    return Array.from(node.children!.values()).map(child => [child.name, child.type])
  }

  async createDirectory(path: string): Promise<void> {
    // 简化版：假设父目录必须存在 (mkdir -p 逻辑较复杂，此处省略)
    const parts = path.split('/').filter(p => p.length > 0)
    const dirName = parts.pop()!
    const parentPath = '/' + parts.join('/')

    const parentNode = this.findNode(parentPath)
    if (!parentNode) throw FileSystemError.FileNotFound(parentPath)
    if (parentNode.children!.has(dirName)) throw FileSystemError.FileExists(path)

    parentNode.children!.set(dirName, {
      type: FileType.Directory,
      name: dirName,
      children: new Map(),
      mtime: Date.now()
    })
  }

  async delete(path: string): Promise<void> {
    const parts = path.split('/').filter(p => p.length > 0)
    const name = parts.pop()!
    const parentPath = '/' + parts.join('/')

    const parentNode = this.findNode(parentPath)
    if (!parentNode || !parentNode.children) throw FileSystemError.FileNotFound(path)

    if (!parentNode.children.has(name)) throw FileSystemError.FileNotFound(path)
    parentNode.children.delete(name)
  }

  async rename(oldPath: string, newPath: string): Promise<void> {
    // 简化版：读取 -> 写入 -> 删除
    const content = await this.readFile(oldPath)
    await this.writeFile(newPath, content)
    await this.delete(oldPath)
  }
}
```

### 3. 核心：VFS 管理器 (挂载与路由)

这是 VFS 的大脑。它维护一张“挂载表”，决定 `/mnt/disk1/file` 应该由哪个 Provider 处理。

```typescript
interface MountPoint {
  path: string
  fs: IFileSystem
  length: number // 路径长度，用于优先匹配最长路径
}

export class VirtualFileSystem implements IFileSystem {
  readonly = false
  private mounts: MountPoint[] = []

  constructor() {
    // 默认挂载一个根内存文件系统，防止 / 无法访问
    this.mount('/', new MemoryFileSystem())
  }

  // --- 挂载机制 ---

  mount(path: string, fs: IFileSystem) {
    // 规范化路径，确保以 / 结尾以便匹配
    const normalizedPath = path.endsWith('/') ? path : path + '/'

    this.mounts.push({
      path: normalizedPath,
      fs: fs,
      length: normalizedPath.length
    })

    // 按路径长度降序排序，确保最长匹配原则
    // 例如：/tmp/logs 应该优先匹配 /tmp/logs 挂载点，而不是 /tmp
    this.mounts.sort((a, b) => b.length - a.length)
    console.log(`[VFS] Mounted ${fs.constructor.name} at ${path}`)
  }

  // --- 路由解析 (Path Resolution) ---

  private resolve(path: string): { fs: IFileSystem; relativePath: string } {
    // 确保路径以 / 开头
    const fullPath = path.startsWith('/') ? path : '/' + path
    // 加上 / 以便匹配挂载点前缀
    const checkPath = fullPath.endsWith('/') ? fullPath : fullPath + '/'

    for (const mount of this.mounts) {
      if (checkPath.startsWith(mount.path)) {
        // 找到了匹配的挂载点
        // 计算相对路径：/mnt/disk1/a.txt (挂载点 /mnt/disk1/) -> /a.txt
        let relativePath = fullPath.slice(mount.path.length - 1)
        if (relativePath === '') relativePath = '/'

        return { fs: mount.fs, relativePath }
      }
    }

    throw new Error(`No filesystem mounted for ${path}`)
  }

  // --- 代理所有操作 ---

  async stat(path: string): Promise<FileStat> {
    const { fs, relativePath } = this.resolve(path)
    return fs.stat(relativePath)
  }

  async readFile(path: string): Promise<Uint8Array> {
    const { fs, relativePath } = this.resolve(path)
    return fs.readFile(relativePath)
  }

  async writeFile(path: string, content: Uint8Array): Promise<void> {
    const { fs, relativePath } = this.resolve(path)
    if (fs.readonly) throw FileSystemError.NoPermissions(path)
    return fs.writeFile(relativePath, content)
  }

  async readDirectory(path: string): Promise<[string, FileType][]> {
    const { fs, relativePath } = this.resolve(path)
    return fs.readDirectory(relativePath)
  }

  async createDirectory(path: string): Promise<void> {
    const { fs, relativePath } = this.resolve(path)
    if (fs.readonly) throw FileSystemError.NoPermissions(path)
    return fs.createDirectory(relativePath)
  }

  async delete(path: string, recursive?: boolean): Promise<void> {
    const { fs, relativePath } = this.resolve(path)
    if (fs.readonly) throw FileSystemError.NoPermissions(path)
    return fs.delete(relativePath, recursive)
  }

  async rename(oldPath: string, newPath: string): Promise<void> {
    const from = this.resolve(oldPath)
    const to = this.resolve(newPath)

    // 关键：跨文件系统移动 (Cross-Device Link)
    if (from.fs !== to.fs) {
      console.warn('[VFS] Cross-device move detected. Falling back to copy-delete.')
      const content = await from.fs.readFile(from.relativePath)
      await to.fs.writeFile(to.relativePath, content)
      await from.fs.delete(from.relativePath)
    } else {
      // 同一文件系统，直接重命名
      return from.fs.rename(from.relativePath, to.relativePath)
    }
  }
}
```

### 4. 扩展：只读网络文件系统 (NetworkFS)

为了演示 VFS 的强大，我们实现一个模拟的只读网络文件系统。

```typescript
export class ReadonlyNetworkFS implements IFileSystem {
  readonly = true
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  async readFile(path: string): Promise<Uint8Array> {
    console.log(`[NetworkFS] Fetching ${this.baseUrl}${path}`)
    // 模拟网络请求
    if (path === '/readme.md') {
      return new TextEncoder().encode('# Remote Readme')
    }
    throw FileSystemError.FileNotFound(path)
  }

  async stat(path: string): Promise<FileStat> {
    if (path === '/readme.md') {
      return { type: FileType.File, size: 100, ctime: 0, mtime: 0 }
    }
    throw FileSystemError.FileNotFound(path)
  }

  // 其他写操作全部抛错
  async writeFile() {
    throw new Error('Readonly FS')
  }
  async delete() {
    throw new Error('Readonly FS')
  }
  async createDirectory() {
    throw new Error('Readonly FS')
  }
  async rename() {
    throw new Error('Readonly FS')
  }
  async readDirectory() {
    return []
  }
}
```

### 5. 实战演示

```typescript
async function main() {
  const vfs = new VirtualFileSystem()

  // 1. 准备不同的后端
  const memFS = new MemoryFileSystem()
  const netFS = new ReadonlyNetworkFS('https://api.example.com/files')

  // 2. 挂载
  // 根目录使用内存 FS
  // /remote 目录使用网络 FS
  vfs.mount('/', memFS)
  vfs.mount('/remote', netFS)

  const encoder = new TextEncoder()
  const decoder = new TextDecoder()

  try {
    // 3. 操作根目录 (MemoryFS)
    console.log('--- Local Operations ---')
    await vfs.writeFile('/local.txt', encoder.encode('Hello Local'))
    const localContent = await vfs.readFile('/local.txt')
    console.log('Read /local.txt:', decoder.decode(localContent))

    // 4. 操作远程目录 (NetworkFS)
    console.log('\n--- Remote Operations ---')
    const remoteContent = await vfs.readFile('/remote/readme.md')
    // 内部路由逻辑：
    // 1. 匹配到 /remote 挂载点
    // 2. 相对路径变为 /readme.md
    // 3. 调用 netFS.readFile('/readme.md')
    console.log('Read /remote/readme.md:', decoder.decode(remoteContent))

    // 5. 跨文件系统移动 (Cross-Device Move)
    console.log('\n--- Cross Device Move ---')
    // 把远程文件“移动”到本地 (实际上是 下载 -> 写入 -> 删除失败)
    // 注意：因为 NetFS 是只读的，delete 会报错，这里仅演示逻辑
    try {
      await vfs.rename('/remote/readme.md', '/local_copy.md')
    } catch (e) {
      console.log('Move failed as expected (Readonly):', e.message)
      // 手动复制
      await vfs.writeFile('/local_copy.md', remoteContent)
      console.log('Manual copy success.')
    }

    // 验证本地新文件
    const copyStat = await vfs.stat('/local_copy.md')
    console.log('Local copy stat:', copyStat)
  } catch (e) {
    console.error('VFS Error:', e)
  }
}

main()
```

### 核心设计思想总结

1.  **抽象层 (Abstraction Layer)**：
    `IFileSystem` 接口是核心。只要实现了这个接口，无论是 IndexedDB、GitHub API 还是 AWS S3，都可以无缝接入这个系统。

2.  **挂载表与最长前缀匹配 (Mount Table & Longest Prefix Match)**：
    这是操作系统文件管理的核心算法。
    如果有 `/` 和 usr 两个挂载点。访问 `/usr/bin/node` 时，必须优先匹配 usr 而不是 `/`。我们在 `mount` 方法中通过按长度降序排序解决了这个问题。

3.  **跨设备操作 (Cross-Device Operations)**：
    在 `rename` 方法中，我们检测了源路径和目标路径是否属于同一个 FS 实例。如果不是，操作系统内核通常会报错（`EXDEV`），但我们在用户态 VFS 中可以做一个“智能回退”：自动执行 `Copy + Delete`，这对上层应用是透明的。

4.  **错误标准化**：
    `FileSystemError` 统一了不同后端的错误码（ENOENT, EACCES 等），这对于上层应用处理异常至关重要。
