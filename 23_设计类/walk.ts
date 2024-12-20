import * as fs from 'fs'
import * as path from 'path'

interface IWalkOptions {
  /**
   * 是否自顶向下遍历.
   *
   * 当 `topdown = true` 时，先遍历当前目录，再遍历子目录.
   * 当 `topdown = false` 时，先遍历子目录，再遍历当前目录.
   *
   * @default true
   */
  topdown?: boolean

  /**
   * 遇到错误时的处理函数.
   *
   * 可以报告错误以继续遍历，或者引发异常以终止遍历.
   */
  onerror?: (err: Error, fullPath: string) => void

  /**
   * 是否跟随符号链接.
   *
   * @default false
   */
  followLinks?: boolean
}

/**
 * 列表中的名称仅仅是名称，不包含路径.
 * 要获取完整路径，需要使用 `path.join(dirPath, name)` 来拼接.
 */
interface IWalkResult {
  dirPath: string

  /**
   * 当 `topdown = true` 时，调用者可以就地修改 `directories` 列表，以控制遍历的子目录(裁剪搜索，或强制指定访问的顺序).
   * 当 `topdown = false` 时，`directories` 列表的修改不会影响遍历的子目录.
   */
  directories: string[]

  files: string[]
}

/**
 * 类似 Python 的 `os.walk` 函数，返回一个生成器，生成 dirpath, dirnames, filenames 三元组.
 *
 * @param top 要遍历的目录
 * @param options 遍历选项 {@link IWalkOptions}
 */
export function* walkSync(
  top: string,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  { topdown = true, onerror = () => {}, followLinks = false }: IWalkOptions = {}
): Generator<IWalkResult, void, unknown> {
  yield* _walkSync(top, topdown, onerror, followLinks)
}

function* _walkSync(
  top: string,
  topdown: boolean,
  onerror: (err: Error, fullPath: string) => void,
  followLinks: boolean
): Generator<IWalkResult, void, unknown> {
  let entries: string[]
  try {
    entries = fs.readdirSync(top)
  } catch (err) {
    onerror(err as Error, top)
    return
  }

  let dirs: string[] = []
  let nondirs: string[] = []
  for (const entry of entries) {
    const fullPath = path.join(top, entry)
    let stat: fs.Stats
    try {
      stat = followLinks ? fs.statSync(fullPath) : fs.lstatSync(fullPath)
    } catch (err) {
      onerror(err as Error, fullPath)
      continue
    }

    if (stat.isDirectory()) {
      dirs.push(entry)
    } else {
      nondirs.push(entry)
    }
  }

  if (topdown) {
    yield { dirPath: top, directories: dirs, files: nondirs }
    for (const d of dirs) {
      yield* _walkSync(path.join(top, d), topdown, onerror, followLinks)
    }
  } else {
    for (const d of dirs) {
      yield* _walkSync(path.join(top, d), topdown, onerror, followLinks)
    }
    yield { dirPath: top, directories: dirs, files: nondirs }
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  for (const result of walkSync('./typings')) {
    console.log(result)
  }
}
