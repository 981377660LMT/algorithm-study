from collections import Counter
from functools import lru_cache
from itertools import permutations

# 请你找出字符串 s 中 重复 k 次的 最长子序列 。
# 如果存在多个满足的子序列，则返回 字典序最大 的那个。如果不存在这样的子序列，返回一个 空 字符串。

# n == s.length
# 2 <= k <= 2000
# 2 <= n < k * 8  (表示不超过7种)
# s 由小写英文字母组成


@lru_cache(None)
def isKRepeatedSubsequence(longer: str, shorter: str, k: int) -> bool:
    """needle是否为pattern的k重复子序列 O(length)"""
    it = iter(longer)
    return all(char in it for char in shorter * k)


class Solution:
    def longestSubsequenceRepeatedK(self, s: str, k: int) -> str:
        counter = Counter(s)
        # 按照数据 不超过7种
        cands = sorted(''.join((key * (counter[key] // k) for key in counter)), reverse=True)
        # 暴力检验  倒序可以提前返回最大字典序
        for length in range(len(cands), 0, -1):
            for perm in permutations(cands, length):
                cur = ''.join(perm)
                if isKRepeatedSubsequence(s, cur, k):
                    return cur
        return ''


for p in permutations('cba', 2):
    print(p)
print(Solution().longestSubsequenceRepeatedK(s="letsleetcode", k=2))
# 输出："let"
# 解释：存在两个最长子序列重复 2 次：let" 和 "ete" 。
# "let" 是其中字典序最大的一个。
