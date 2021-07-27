String.prototype.longestPalindrome = function () {
  var str = '#', //预处理字符串
    mid = 0, //当前最长回文子串的中心
    right = 0, //当前最长回文子串的右边界
    maxLen = 0, //最长回文子串长度
    maxLenMid = 0, //最长回文子串的中心
    child = [] //存放每个字符的回文长度
  /*
      预处理字符串,之后的字符串长度奇偶个数就统一了
      比如  abba   -->  #a#b#b#a#     偶数个长度变为了奇数个
            abfba  -->  #a#b#f#b#a#   奇数个长度没有变
    */
  for (var i = 0; i < this.length; i++) {
    str += this[i] + '#'
  }
  for (var i = 0; i < str.length; i++) {
    /*第 i 个字符是否还在右边界内
        Math.min(child[2*mid-i],right-i) 解释:这里为什么选两者中小的那个呢？
        1、当child[2*mid-i]较小时，说明以 i 点处字符为中心的回文*完全*在(mid,right)区间内，直接附值即可。
        2、当child[2*mid-i]较大时，说明以 i 点处字符为中心的回文*至少*在(mid,right)区间内，可能有
           超出的能够匹配的字符，所以先附值 right-i,然后再对超出 right 边界的字符一一做对称匹配。 
        算出 child[i] 的值了，下面的 while 循环就是对接下来的字符做匹配操作的
      */
    child[i] = i < right ? Math.min(child[2 * mid - i], right - i) : 1
    //接下来从 i+-child[i] 再做逐个比较
    while (str.charAt(i + child[i]) == str.charAt(i - child[i])) {
      child[i]++
    }
    //是否更新右边界:如果当前 i 加上 i 的回文长度 child[i] 大于原先的右边界,就更新
    if (right < child[i] + i) {
      mid = i
      right = child[i] + i
    }
    //是否更新最长回文子串的长度及其中心位置：如果当前 i 的回文长度child[i]大于原先的最长回文长度,就更新
    if (maxLen < child[i]) {
      maxLen = child[i]
      maxLenMid = i
    }
  }
  //根据两个变量 maxLenMid 和 maxLen 判断子串在原字符串中的位置,取出子串再返回
  return this.substr((maxLenMid + 1 - maxLen) / 2, maxLen - 1)
}

console.log('asas'.longestPalindrome())
