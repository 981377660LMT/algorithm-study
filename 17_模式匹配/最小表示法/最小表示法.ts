// 字符串最小表示法 O（n)算法
// 有一个首位相连的字符串，我们要寻找一个位置，从这个位置向后形成一个新字符串，我们需要使这个字符串字典序最小。

// 我们这里要i = 0,j = 1,k = 0,表示从i开始k长度和从j开始k长度的字符串相同（i，j表示当前判断的位置）
// 当我们str[i] == str[j]时，根据上面k的定义，我们的需要进行k+1操作
// 当str[i] > str[j]时，我们发现i位置比j位置上字典序要大，那么不能使用i作为开头了，我们要将i向后移动，移动多少呢？有因为i开头和j开头的有k个相同的字符，那么就执行 i + = k +1
// 相反str[i] < str[j]时，执行：j = j + k +1
// 滑动方式有个小细节，若滑动后i == j，将正在变化的那个指针再+1
// 最终i和j中较小的值就是我们最终开始的位置
// 相反如果是最大表示法的话，我们就要求解字典序最大的字符串，那么我们只需要在执行第二或第三个操作时选择较大的那个位置较好了
// https://a-wimpy-boy.blog.csdn.net/article/details/80136776?spm=1001.2101.3001.6650.1&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-1.no_search_link&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7Edefault-1.no_search_link
function findMinimunIsomorphic(str: string): string {
  if (str.length <= 1) return str

  const n = str.length
  let i = 0
  let j = 1
  let k = 0

  while (i < n && j < n && k < n) {
    const diff = str.codePointAt((i + k) % n)! - str.codePointAt((j + k) % n)!

    if (diff === 0) {
      k++
      continue
    }

    if (diff > 0) i += k + 1
    else if (diff < 0) j += k + 1

    if (i === j) j++

    k = 0
  }

  const res = i > j ? j : i
  console.log(i, j)
  return str.slice(res) + str.slice(0, res)
}

console.log(findMinimunIsomorphic('bcaijab'))
