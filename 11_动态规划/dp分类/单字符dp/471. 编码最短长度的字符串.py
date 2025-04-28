# 给定一个 非空 字符串，将其编码为具有最短长度的字符串。
# 编码规则是：k[encoded_string]，其中在方括号 encoded_string 中的内容重复 k 次。
# 如果编码的过程不能使字符串缩短，则不要对其进行编码
# 1 <= s.length <= 150

from functools import lru_cache


def minimalPeriod(s: str) -> int:
    """计算字符串 s 的最小周期。如果没有找到，返回 len(s)。"""
    n = len(s)
    res = (s + s).find(s, 1, -1)
    return res if res != -1 else n


class Solution:
    @lru_cache(None)
    def encode(self, s: str) -> str:
        res = s
        if len(s) <= 4:
            return res
        period = minimalPeriod(s)
        if period < len(s):
            res = str(len(s) // period) + "[" + self.encode(s[:period]) + "]"
        for i in range(1, len(s)):
            left = self.encode(s[:i])
            right = self.encode(s[i:])
            res = min(res, left + right, key=len)
        return res


print(Solution().encode(s="aabcaabcd"))
# 输出："2[aabc]d"
# 解释："aabc" 出现两次，因此一种答案可以是 "2[aabc]d"。
