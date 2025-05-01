# 1638. 统计只差一个字符的子串数目
# https://leetcode.cn/problems/count-substrings-that-differ-by-one-character/
# 相似题目
# 795. 区间子数组个数
# 2444. 统计定界子数组的数目
#
# 用三个参数表示一个对子串 (sStart, sEnd, tEnd)


class Solution:
    def countSubstrings(self, s: str, t: str) -> int:
        res, n, m = 0, len(s), len(t)
        for d in range(1 - m, n):  # d=i-j, j=i-d
            i = max(d, 0)
            k0 = k1 = i - 1
            while i < n and i - d < m:
                if s[i] != t[i - d]:
                    k0 = k1  # 上上一个不同
                    k1 = i  # 上一个不同
                res += k1 - k0
                i += 1
        return res
