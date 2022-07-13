#剑指 Offer 03 数组中重复的数字
`原地哈希`

#剑指 Offer 05. 替换空格
`预留空间+倒序双指针`

剑指 Offer 06. 从尾到头打印链表
`后序dfs`

#剑指 Offer 10- I 斐波那契数列
`矩阵快速幂`
![图 2](../images/e4d6d0dd1cf0c785986fecca1fe8c9d51c207b60a1e0a5f1dfadb882916d3a79.png)

#剑指 Offer 13 机器人的运动范围
`dfs(i, j, si, sj)`

#剑指 Offer 14- II 剪绳子 II
`大于4尽量剪3`

#剑指 Offer 17 打印从 1 到最大的 n 位数
`bt(index: number, digit: number, path: string[])`

#剑指 Offer 21 调整数组顺序使奇数位于偶数前面
`头尾双指针`

#剑指 Offer 26 树的子结构
`双递归：return dfs(A, B) || isSubStructure(A.left, B) || isSubStructure(A.right, B)`

#剑指 Offer 31 栈的压入、弹出序列
`模拟,维护popIndex`

#剑指 Offer 33 二叉搜索树的后序遍历序列
`左、右、根=>找根和左右子树的分界点`

#剑指 Offer 36 二叉搜索树与双向链表
`中序遍历，记录pre`

#剑指 Offer 38 字符串的排列
`if(i>0&&str[i]===str[i-1]&&visited[i-1]) continue`

#剑指 Offer 50 第一个只出现一次的字符
`LinkedHashMap`

#剑指 Offer 59 - II 队列的最大值
`用一个deque保存正常元素，另一个deque保存单调递减的元素`

#剑指 Offer 60 n 个骰子的点数
`dp[n][i]`

#剑指 Offer 61 扑克牌中的顺子
`排序，统计大小王`

#剑指 Offer II 003 前 n 个数字二进制中 1 的个数
`res[i] = i % 2 === 0 ? res[i / 2] : res[i - 1] + 1`

#剑指 Offer II 009 乘积小于 K 的子数组
`product;滑窗 `

#剑指 Offer II 014 字符串中的变位词
`记录need以及needCounter`

#剑指 Offer II 020 回文子字符串的个数
`中心扩展`

#剑指 Offer II 034 外星语言是否排序
`母题`
`没有break就比较长度`

#剑指 Offer II 057 值和下标之差都在给定的范围内
`值得差：分桶`
`有序集合:二分判断index左右两侧的数据与num差的绝对值是否小于等于t`

#剑指 Offer II 064. 神奇的字典
`广义邻居：word[:i]+'*'+word[i+1:] => word`
`广义邻居：127. Word Ladder`

#剑指 Offer II 067 最大的异或
`异或前缀树`

#剑指 Offer II 081 允许重复选择元素的组合
`candidates.sort((a, b) => a - b)`
`if (i > index && candidates[i] === candidates[i - 1]) continue`

#剑指 Offer II 082 含有重复元素集合的组合
`candidates.sort((a, b) => a - b)`
`if (i > index && candidates[i] === candidates[i - 1]) continue`

#剑指 Offer II 083 没有重复元素集合的全排列
`bt(path: number[], visited: number)`

#剑指 Offer II 084 含有重复元素集合的全排列
`if (visited & (1 << i)) continue`
`if(i>0&&nums[i]===nums[i-1]&&visited&(1<<(i-1))) continue`

#剑指 Offer II 107 矩阵中的距离
`多源bfs`

#剑指 Offer II 115 重建序列
`拓扑排序唯一：queue始终长度为1`
