/**
 * @param {string} input
 * @return {number}
 * 每个文件名遵循 name.extension
 * 返回文件系统中 指向文件的最长绝对路径 的长度。 如果系统中没有文件，返回 0。
 */
function lengthLongestPath(input: string): number {
  // 用split(‘\n’) 将原串分割开，相当于一次读一行
  // 利用字符串前面’\t’的个数来当前目录/文件在第几层
  const dir = input.split('\n')
  console.log(dir)
  let max = 0
  const pre = Array<number>(dir.length + 1).fill(0)
  for (let i = 0; i < dir.length; i++) {
    const fileName = dir[i]
    const depth = fileName.lastIndexOf('\t') + 2 // 没有\t 代表第一层
    const len = fileName.length - (depth - 1) //长度要去除'\t'

    if (fileName.includes('.')) {
      max = Math.max(max, pre[depth - 1] + len + depth - 1)
    } else {
      pre[depth] = pre[depth - 1] + len
    }
  }

  return max
}

// \n换行\t制表
console.log(lengthLongestPath(`dir\n\tsubdir1\n\tsubdir2\n\t\tfile.ext`))
// 只有一个文件，绝对路径为 "dir/subdir2/file.ext" ，路径长度 20
console.log(lengthLongestPath(`file1.txt\nfile2.txt\nlongfile.txt`))

export {}
console.log('file.ext'.lastIndexOf('\t'))
console.log('\tfile.ext'.lastIndexOf('\t'))
console.log('\t\tfile.ext'.lastIndexOf('\t'))
