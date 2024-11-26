import fs from 'fs'
import path from 'path'

export function walk(rootPath: string, f: (filePath: string) => boolean | void): void {
  const visited = new Set<string>()
  let queue: string[] = [rootPath]

  while (queue.length) {
    const nextQueue: string[] = []

    for (const curPath of queue) {
      const realPath = fs.realpathSync(curPath)
      if (visited.has(realPath)) continue
      visited.add(realPath)

      if (!fs.statSync(curPath).isDirectory()) {
        if (f(curPath)) return
      } else {
        for (const child of fs.readdirSync(curPath)) {
          const childPath = path.join(curPath, child)
          nextQueue.push(childPath)
        }
      }
    }

    queue = nextQueue
  }
}

export function findLowestFile(rootPath: string, predicate: (filePath: string) => boolean): string | undefined {
  let res: string | undefined
  walk(rootPath, filePath => {
    if (predicate(filePath)) {
      res = filePath
      return true
    }
  })
  return res
}
