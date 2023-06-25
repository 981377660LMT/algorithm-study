from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你三个整数 x ，y 和 z 。

# 这三个整数表示你有 x 个 "AA" 字符串，y 个 "BB" 字符串，和 z 个 "AB" 字符串。你需要选择这些字符串中的部分字符串（可以全部选择也可以一个都不选择），将它们按顺序连接得到一个新的字符串。新字符串不能包含子字符串 "AAA" 或者 "BBB" 。

# 请你返回新字符串的最大可能长度。


# 子字符串 是一个字符串中一段连续 非空 的字符序列。
class Solution:
    def longestString(self, x: int, y: int, z: int) -> int:
        if x == y:
            return 2 * (x + y + z)
        min_ = min(x, y)
        return (min_ + min_ + 1 + z) * 2


# 1
# 39
# 14
# 预期:34


print(Solution().longestString(1, 39, 14))
