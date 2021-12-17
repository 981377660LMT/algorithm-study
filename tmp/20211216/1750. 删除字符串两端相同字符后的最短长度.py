# 输入：s = "aabccabba"
# 输出：3
# 解释：最优操作序列为：
# - 选择前缀 "aa" 和后缀 "a" 并删除它们，得到 s = "bccabb" 。
# - 选择前缀 "b" 和后缀 "bb" 并删除它们，得到 s = "cca" 。


class Solution:
    def minimumLength(self, s: str) -> int:
        left, right, cur = 0, len(s) - 1, ''
        while left < right and s[left] == s[right]:
            cur = s[left]
            while left < right and s[left] == cur:
                left += 1
            while left < right and s[right] == cur:
                right -= 1

        return 0 if s[left] == cur else right - left + 1

    # O(n) * O(string generation)
    def minimumLength2(self, s: str) -> int:
        while len(s) > 1 and s[0] == s[-1]:
            s = s.strip(s[0])
        return len(s)

