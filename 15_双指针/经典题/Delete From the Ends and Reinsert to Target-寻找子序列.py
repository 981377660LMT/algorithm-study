# s/t互为anagram
# 每次都可以将s的头或尾插入到t的任意一个位置
# 求最小转换次数

# n ≤ 1,000
class Solution:
    def solve(self, s, t):
        res = 0

        # 从每s串的个位置开始寻找相同子序列的长度(不改变相邻位置) 注意不是LCS
        for i in range(len(s)):
            j, k = 0, i
            while j < len(s) and k < len(s):
                if s[k] == t[j]:
                    k += 1
                j += 1
            res = max(res, k - i)

        return len(s) - res


# class Solution:
#     def solve(self, s, t):
#         @lru_cache(None)
#         def lcs(a, b):
#             if not a or not b:
#                 return 0
#             for i in range(len(b)):
#                 if a[0] == b[i]:
#                     return 1 + lcs(a[1:], b[i + 1 :])
#             return 0

#         return len(s) - max(lcs(s[i:], t) for i in range(len(s)))
