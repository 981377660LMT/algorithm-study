后缀数组（Suffix Array）主要是两个数组：SA 和 RK
SA[i] 表示将所有后缀(n 个)排序后第 i 小的后缀是第几个后缀
RK[i] 表示第 i 个后缀 的排名
Height[i] 来记录排名为 i 的非空后缀与排名为 i−1 的非空后缀的最长公共前缀的长度
![](image/note/1651156168474.png)

应用

- 从字符串首尾取字符最小化字典序
- height 数组
  height 数组的定义
  H[i]=LCP(SA[i],SA[i-1])
  即`第 i 名的后缀与它前一名的后缀的最长公共前缀`

1. 后缀数组处理两个串的技巧:**串起来**
   如果要比较 S 和 T 的所有子串的字典序
   那么可以令新的字符串为 `S + '#' + T + '|'`
   其中 '#' 是字典序很小的字符, '|' 是字典序很大的字符
   这样可以保证`两个字符串在比较完长度为 n 后 ,S 后面的#小于 T 中任意一个字符`
   然后对这个新的字符串求后缀数组 从而得到原来的 S 和 T 的所有子串的排名

   ```go
   // 例子：寻找最长公共子串
   dummy := max(maxs(ords1...), maxs(ords2...)) + 1
   sb := make([]int, 0, len(ords1)+len(ords2)+1)
   sb = append(sb, ords1...)
   sb = append(sb, dummy)
   sb = append(sb, ords2...)
   sa, _, lcp := UseSA(sb)

   len_ := 0
   for i := 1; i < len(sb); i++ {
   	if (sa[i-1] < len(ords1)) == (sa[i] < len(ords1)) {  // 来自同一个串
   		continue
   	}
   	if lcp[i] <= len_ {  // 与上一个公共子串长度相同或更短
   		continue
   	}
   	len_ = lcp[i]

   	// 来自s和t的不同子串
   	// 找到了(严格)更长的公共子串,更新答案
   	i1, i2 := sa[i-1], sa[i]
   	if i1 > i2 {
   		i1, i2 = i2, i1
   	}

   	start1 = i1
   	end1 = start1 + len_
   	start2 = i2 - len(ords1) - 1
   	end2 = start2 + len_
   }
   ```
