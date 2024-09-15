import fs from 'fs'
import path from 'path'

interface IWalkOptions {}

export function walk(rootPath: string, f: (filePath: string) => boolean | void, options: IWalkOptions = {}): void {
  let queue: string[] = [rootPath]
  while (queue.length) {
    const nextQueue: string[] = []
    for (const curPath of queue) {
      if (!fs.statSync(curPath).isDirectory()) {
        if (f(curPath)) return
      } else {
        for (const child of fs.readdirSync(curPath)) {
          nextQueue.push(path.join(curPath, child))
        }
      }
    }
    queue = nextQueue
  }
}

export function findLowestFile(rootPath: string, predicate: (filePath: string) => boolean): string | undefined {
  let res: string | undefined = undefined
  walk(rootPath, filePath => {
    if (predicate(filePath)) {
      res = filePath
      return true
    }
  })
  return res
}
