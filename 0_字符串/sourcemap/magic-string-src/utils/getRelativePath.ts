/**
 * 计算从一个文件路径到另一个文件路径的相对路径。
 *
 * 根据路径前缀定位到lca，然后分别计算两个路径到lca的相对路径，然后拼接。
 */
export default function getRelativePath(from: string, to: string) {
  const fromPath = from.split(/[/\\]/)
  const toPath = to.split(/[/\\]/)

  // fromPath.pop() // !get dirname

  {
    // remove common prefix
    let lcaDepth = 0
    while (fromPath[lcaDepth] === toPath[lcaDepth]) lcaDepth++
    fromPath.splice(0, lcaDepth)
    toPath.splice(0, lcaDepth)
  }

  fromPath.fill('..')

  return fromPath.concat(toPath).join('/')
}

if (require.main === module) {
  console.log(getRelativePath('a/b/c', 'a/b/c/d/e')) // d/e

  const from = '/a/b/c/file.txt'
  const to = '/a/d/e/target.txt'
  console.log(getRelativePath(from, to)) // ../../../d/e/target.txt
}
