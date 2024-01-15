from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的字符串 s 、字符串 a 、字符串 b 和一个整数 k 。

# 如果下标 i 满足以下条件，则认为它是一个 美丽下标：


# 0 <= i <= s.length - a.length
# s[i..(i + a.length - 1)] == a
# 存在下标 j 使得：
# 0 <= j <= s.length - b.length
# s[j..(j + b.length - 1)] == b
# |j - i| <= k
# 以数组形式按 从小到大排序 返回美丽下标。
class Solution:
    def beautifulIndices(self, s: str, a: str, b: str, k: int) -> List[int]:
        ...
