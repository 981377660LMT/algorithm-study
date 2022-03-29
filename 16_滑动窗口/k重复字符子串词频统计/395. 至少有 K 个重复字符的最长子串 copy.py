class Solution:
    def longestSubstring(self, s: str, k: int) -> int:
        if not s:
            return 0
        for char in set(s):
            if s.count(char) < k:
                return max(self.longestSubstring(t, k) for t in s.split(char))
        return len(s)
