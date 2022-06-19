# 1 <= s.length <= 1000
# 1 <= k <= 1e9


class Solution:
    def longestSubsequence(self, s: str, k: int) -> int:
        """s 的 最长 子序列，且该子序列对应的 二进制 数字小于等于 k 
        
        最长子序列 包含 原二进制字符串 中的所有的 0
        倒序判断每个数选不选 优先保留较低位的 1
        """

        res = ''
        for i in range(len(s) - 1, -1, -1):
            cand = s[i] + res
            if int(cand, 2) <= k:
                res = cand
        return len(res)


# !最长二进制子序列一定是由小于等于k的最长后缀加上前导0数量。

# print(Solution().longestSubsequence(s="1001010", k=5))
# print(Solution().longestSubsequence(s="1011", k=281854076))
# print(Solution().longestSubsequence(s="101", k=281854076))
# print(Solution().longestSubsequence(s="00101001", k=1))
# 5 4
print(
    Solution().longestSubsequence(
        s="1111111111111111111111111111111111111111111111111111111111111111", k=1000000000
    )
)
# 29
