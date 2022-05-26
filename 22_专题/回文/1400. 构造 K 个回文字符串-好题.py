from collections import Counter

# 如果你可以用 s 中所有字符构造 k 个回文字符串，那么请你返回 True

# 奇数的频率的字母种数不能超过k
class Solution:
    def canConstruct(self, s: str, k: int) -> bool:
        if k > len(s):
            return False
        counter = Counter(s)
        return sum(c & 1 for c in counter.values()) <= k


print(Solution().canConstruct(s="annabelle", k=2))
# 输出：true
# 解释：可以用 s 中所有字符构造 2 个回文字符串。
# 一些可行的构造方案包括："anna" + "elble"，"anbna" + "elle"，"anellena" + "b"

