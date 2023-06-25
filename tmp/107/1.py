from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的数组 words ，数组中包含 互不相同 的字符串。

# 如果字符串 words[i] 与字符串 words[j] 满足以下条件，我们称它们可以匹配：

# 字符串 words[i] 等于 words[j] 的反转字符串。
# 0 <= i < j < words.length
# 请你返回数组 words 中的 最大 匹配数目。


# 注意，每个字符串最多匹配一次。
class Solution:
    def maximumNumberOfStringPairs(self, words: List[str]) -> int:
        counter = Counter(words)
        res = 0
        for k, v in counter.items():
            need = k[::-1]
            if need in counter and need != k:
                res += 1
        return res // 2


# ["ff","tx","qr","zw","wr","jr","zt","jk","sq","xx"]
print(
    Solution().maximumNumberOfStringPairs(
        ["ff", "tx", "qr", "zw", "wr", "jr", "zt", "jk", "sq", "xx"]
    )
)
