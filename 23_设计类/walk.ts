import * as fs from 'fs'
import * as path from 'path'

type WalkCallback = (dirpath: string, dirnames: string[], filenames: string[]) => void

interface IWalkOptions {
  topdown?: boolean
  onerror?: (error: NodeJS.ErrnoException) => void
  followlinks?: boolean
}

export function walk(top: string, callback: WalkCallback, options: IWalkOptions = {}): void {
  const {
    topdown = true,
    onerror = (err: NodeJS.ErrnoException) => {
      throw err
    },
    followlinks = false
  } = options

  const queue: string[] = [top]

  while (queue.length > 0) {
    const currentDir = queue.shift() as string

    let dirents: fs.Dirent[]
    try {
      dirents = fs.readdirSync(currentDir, { withFileTypes: true })
    } catch (error) {
      onerror(error as NodeJS.ErrnoException)
      continue
    }

    const dirnames: string[] = []
    const filenames: string[] = []

    for (const dirent of dirents) {
      if (dirent.isDirectory()) {
        const fullPath = path.join(currentDir, dirent.name)
        dirnames.push(dirent.name)
        if (followlinks) {
          try {
            const stats = fs.lstatSync(fullPath)
            if (stats.isSymbolicLink()) {
              const realPath = fs.realpathSync(fullPath)
              queue.push(realPath)
              continue
            }
          } catch (error) {
            onerror(error as NodeJS.ErrnoException)
            continue
          }
        }
        queue.push(path.join(currentDir, dirent.name))
      } else {
        filenames.push(dirent.name)
      }
    }

    if (topdown) {
      callback(currentDir, dirnames, filenames)
    }

    if (!topdown) {
      callback(currentDir, dirnames, filenames)
    }
  }
}
