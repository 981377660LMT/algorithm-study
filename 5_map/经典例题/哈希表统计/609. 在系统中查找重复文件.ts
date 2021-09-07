/**
 * @param {string[]} paths
 * @return {string[][]}
 */
const findDuplicate = function (paths: string[]): string[][] {
  const map = new Map<string, string[]>()

  for (const path of paths) {
    const [root, ...files] = path.split(/\s+/)
    for (const file of files) {
      const [fileName, content] = file.split('(')
      !map.has(content) && map.set(content, [])
      map.get(content)!.push(`${root}/${fileName}`)
    }
  }

  return [...map.values()].filter(v => v.length >= 2)
}

console.log(
  findDuplicate([
    'root/a 1.txt(abcd) 2.txt(efgh)',
    'root/c 3.txt(abcd)',
    'root/c/d 4.txt(efgh)',
    'root 4.txt(efgh)',
  ])
)
// 输出：
// [["root/a/2.txt","root/c/d/4.txt","root/4.txt"],["root/a/1.txt","root/c/3.txt"]]

export default 1
