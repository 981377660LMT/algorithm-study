from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个正整数 low 和 high 。

# 对于一个由 2 * n 位数字组成的整数 x ，如果其前 n 位数字之和与后 n 位数字之和相等，则认为这个数字是一个对称整数。


# 返回在 [low, high] 范围内的 对称整数的数目 。
class Solution:
    def countSymmetricIntegers(self, low: int, high: int) -> int:
        res = 0
        for v in range(low, high + 1):
            s = str(v)
            if len(s) % 2 == 0:
                sum1 = sum([int(x) for x in s[: len(s) // 2]])
                sum2 = sum([int(x) for x in s[len(s) // 2 :]])
                if sum1 == sum2:
                    res += 1
        return res
