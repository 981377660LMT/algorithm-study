# !堵车
# 3228. 将 1 移动到末尾的最大操作次数
# https://leetcode.cn/problems/maximum-number-of-operations-to-move-ones-to-the-end/solutions/2851730/du-che-pythonjavacgo-by-endlesscheng-tllv/
# 给你一个 二进制字符串 s。
# 你可以对这个字符串执行 任意次 下述操作：
# 选择字符串中的任一下标 i（ i + 1 < s.length ），该下标满足 s[i] == '1' 且 s[i + 1] == '0'。
# 将字符 s[i] 向 右移 直到它到达字符串的末端或另一个 '1'。例如，对于 s = "010010"，如果我们选择 i = 1，结果字符串将会是 s = "000110"。
# 返回你能执行的 最大 操作次数。
#
# 把 1 当作车，想象有一条长为 n 的道路上有一些车。
# 我们要把所有的车都排到最右边，例如 011010 最终要变成 000111。
# 如果我们优先操作右边的车，那么每辆车都只需操作一次
# !如果我们优先操作左边的（能移动的）车，这会制造大量的「堵车」，那么每辆车的操作次数就会更多


from itertools import groupby


class Solution:
    def maxOperations(self, s: str) -> int:
        c1 = 0
        res = 0
        for kind, group in groupby(s):
            len_ = len(list(group))
            if kind == "1":
                c1 += len_
            else:
                res += c1
        return res
