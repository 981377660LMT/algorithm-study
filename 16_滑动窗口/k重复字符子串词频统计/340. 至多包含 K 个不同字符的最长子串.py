from collections import defaultdict


class Solution:
    def lengthOfLongestSubstringKDistinct(self, s: str, k: int) -> int:
        """至多包含 K 个不同字符的最长子串"""
        res, left, n = 0, 0, len(s)
        counter = defaultdict(int)
        for right in range(n):
            counter[s[right]] += 1
            while left <= right and len(counter) > k:
                counter[s[left]] -= 1
                if counter[s[left]] == 0:
                    del counter[s[left]]
                left += 1
            res = max(res, right - left + 1)
        return res
