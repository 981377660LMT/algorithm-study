# 每次right变化
# !看left 最多到哪里，left这之前的都是可以的
class Solution:
    def numberOfSubstrings(self, s: str) -> int:
        res = left = 0
        counter = {char: 0 for char in "abc"}
        for right in range(len(s)):
            counter[s[right]] += 1
            while left <= right and all(counter.values()):
                counter[s[left]] -= 1
                left += 1
            res += left
        return res


print(Solution().numberOfSubstrings(s="abcabc"))
# 输出：10
# 解释：包含 a，b 和 c 各至少一次的子字符串为 "abc", "abca", "abcab", "abcabc", "bca", "bcab", "bcabc", "cab", "cabc" 和 "abc" (相同字符串算多次)。
