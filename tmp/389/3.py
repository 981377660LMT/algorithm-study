from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 word 和一个整数 k。

# 如果 |freq(word[i]) - freq(word[j])| <= k 对于字符串中所有下标 i 和 j  都成立，则认为 word 是 k 特殊字符串。

# 此处，freq(x) 表示字符 x 在 word 中的出现频率，而 |y| 表示 y 的绝对值。


# 返回使 word 成为 k 特殊字符串 需要删除的字符的最小数量。
class Solution:
    def minimumDeletions(self, word: str, k: int) -> int:
        counter = Counter(word)
        freq = sorted(counter.values())

        # 枚举最小值
        res = len(word)
        for min_ in freq:
            max_ = min_ + k
            curCost = 0
            for f in freq:
                if f < min_:
                    curCost += f
                if f > max_:
                    curCost += f - max_
            res = min(res, curCost)

        return res


# word = "aabcaba", k = 0
print(Solution().minimumDeletions("aabcaba", 0))
# word = "aaabaaa", k = 2
print(Solution().minimumDeletions("aaabaaa", 2))
