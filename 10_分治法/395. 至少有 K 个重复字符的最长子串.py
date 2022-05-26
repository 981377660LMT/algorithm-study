# 分治，代码简洁
# O(N∗26∗26)
class Solution:
    def longestSubstring(self, s: str, k: int) -> int:
        if not s:
            return 0
        for char in set(s):
            if s.count(char) < k:
                return max(self.longestSubstring(t, k) for t in s.split(char))
        return len(s)


print(Solution().longestSubstring("ababbc", 2))
# 输出：5
# 解释：最长子串为 "ababb" ，其中 'a' 重复了 2 次， 'b' 重复了 3 次。
