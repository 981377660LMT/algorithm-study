from itertools import groupby
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 word，请你使用以下算法进行压缩：


# 从空字符串 comp 开始。当 word 不为空 时，执行以下操作：
# 移除 word 的最长单字符前缀，该前缀由单一字符 c 重复多次组成，且该前缀长度 最多 为 9 。
# 将前缀的长度和字符 c 追加到 comp 。
# 返回字符串 comp 。
class Solution:
    def compressedString(self, word: str) -> str:
        groups = [(char, len(list(group))) for char, group in groupby(word)]
        res = []
        for char, cnt in groups:
            div, mod = divmod(cnt, 9)
            for _ in range(div):
                res.append(f"9{char}")
            if mod:
                res.append(f"{mod}{char}")
        return "".join(res)
