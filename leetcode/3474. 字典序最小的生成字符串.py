# 3474. 字典序最小的生成字符串
# https://leetcode.cn/problems/lexicographically-smallest-generated-string/solutions/3591754/gou-zao-tan-xin-by-tsreaper-4gtj/
#
# 给你两个字符串，str1 和 str2，其长度分别为 n 和 m 。
# 如果一个长度为 n + m - 1 的字符串 word 的每个下标 0 <= i <= n - 1 都满足以下条件，则称其由 str1 和 str2 生成：
# 如果 str1[i] == 'T'，则长度为 m 的 子字符串（从下标 i 开始）与 str2 相等，即 word[i..(i + m - 1)] == str2。
# 如果 str1[i] == 'F'，则长度为 m 的 子字符串（从下标 i 开始）与 str2 不相等，即 word[i..(i + m - 1)] != str2。
# 返回可以由 str1 和 str2 生成 的 字典序最小 的字符串。如果不存在满足条件的字符串，返回空字符串 ""。
# n<=1e4,m<=500
#
# 构造+贪心 O(n*m)
#
# 先填T，再填F，不对就回退
# Z函数线性解法 https://leetcode.cn/problems/lexicographically-smallest-generated-string/solutions/3592764/liang-chong-fang-fa-tan-xin-bao-li-pi-pe-gxyn/


class Solution:
    def generateString(self, str1: str, str2: str) -> str:
        n, m = len(str1), len(str2)
        flag = [False] * (n + m - 1)  # flag[i]表示下标i是否已经填充
        sb = [""] * (n + m - 1)

        # 先填T
        for i, b in enumerate(str1):
            if b == "T":
                for j in range(m):
                    if flag[i + j] and sb[i + j] != str2[j]:
                        return ""
                    flag[i + j] = True
                    sb[i + j] = str2[j]

        # 没填充的，全部用a填充
        for i in range(n + m - 1):
            if not flag[i]:
                sb[i] = "a"

        # 再填F
        for i, b in enumerate(str1):
            if b == "F" and "".join(sb[i : i + m]) == str2:
                failed = True
                for j in range(m - 1, -1, -1):
                    if not flag[i + j]:
                        sb[i + j] = "b"
                        failed = False
                        break
                if failed:
                    return ""

        return "".join(sb)
