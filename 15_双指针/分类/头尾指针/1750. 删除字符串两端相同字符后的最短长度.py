"""删除字符串两端相同字符后的最短长度"""


class Solution:
    def minimumLength(self, s: str) -> int:
        left, right = 0, len(s) - 1
        while left < right and s[left] == s[right]:
            cur = s[left]
            while left <= right and s[left] == cur:
                left += 1
            while left <= right and s[right] == cur:
                right -= 1
        return right - left + 1


print(Solution().minimumLength(s="aabccabba"))
# 输出：3
# 解释：最优操作序列为：
# - 选择前缀 "aa" 和后缀 "a" 并删除它们，得到 s = "bccabb" 。
# - 选择前缀 "b" 和后缀 "bb" 并删除它们，得到 s = "cca" 。
