from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 num 。你知道 Danny Mittal 会偷偷将 0 到 9 中的一个数字 替换 成另一个数字。

# 请你返回将 num 中 恰好一个 数字进行替换后，得到的最大值和最小值的差位多少。

# 注意：

# 当 Danny 将一个数字 d1 替换成另一个数字 d2 时，Danny 需要将 nums 中所有 d1 都替换成 d2 。
# Danny 可以将一个数字替换成它自己，也就是说 num 可以不变。
# Danny 可以将数字分别替换成两个不同的数字分别得到最大值和最小值。
# 替换后得到的数字可以包含前导 0 。
# Danny Mittal 获得周赛 326 前 10 名，让我们恭喜他。
class Solution:
    def minMaxDifference(self, num: int) -> int:
        chars = list(str(num))
        mp = {i: i for i in range(10)}
        res = []
        for fromNum in range(10):
            for toNum in range(10):
                mp[fromNum] = toNum
                newNum = int("".join([str(mp[int(c)]) for c in chars]))
                res.append(newNum)
                mp[fromNum] = fromNum
        return max(res) - min(res)
