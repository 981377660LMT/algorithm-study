// 给定两个字符串，编写一个函数判定它们是否只需要一次(或者零次)编辑。

// 如果长度之差大于一的话，显然不可能一次编辑。
// 如果长度之差满足要求，那么只需要从到遍历两字符串，找到第一个不相同的字符所在的位置i。
// 对于i位置字符，first到second编辑方法又有三种可能：first删除一个，增加一个，修改一个那么就对应了代码中first[i:] == second[i+1:] or first[i+1:] == second[i+1:] or first[i+1:] == second[i:]的三种可能，满足一个即可。

// 作者：xiao-zhu-ssp
// 链接：https://leetcode-cn.com/problems/one-away-lcci/solution/mian-shi-ti-05yi-ci-bian-ji-liang-jie-by-qzue/
// 来源：力扣（LeetCode）
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

function oneEditAway(first: string, second: string): boolean {
  if (Math.abs(first.length - second.length) >= 2) return false
  const len = Math.min(first.length, second.length)
  for (let i = 0; i < len; i++) {
    if (first[i] === second[i]) continue
    // 增，删，改
    return (
      first.slice(i + 1) === second.slice(i) ||
      first.slice(i) === second.slice(i + 1) ||
      first.slice(i + 1) === second.slice(i + 1)
    )
  }
  return true
}

console.log(oneEditAway('pale', 'ple'))
