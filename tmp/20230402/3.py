from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 有两只老鼠和 n 块不同类型的奶酪，每块奶酪都只能被其中一只老鼠吃掉。

# 下标为 i 处的奶酪被吃掉的得分为：

# 如果第一只老鼠吃掉，则得分为 reward1[i] 。
# 如果第二只老鼠吃掉，则得分为 reward2[i] 。
# 给你一个正整数数组 reward1 ，一个正整数数组 reward2 ，和一个非负整数 k 。

# 请你返回第一只老鼠恰好吃掉 k 块奶酪的情况下，最大 得分为多少。
class Solution:
    def miceAndCheese(self, reward1: List[int], reward2: List[int], k: int) -> int:
        R = [(a, a - b, i) for i, (a, b) in enumerate(zip(reward1, reward2))]
        R.sort(key=lambda x: x[1], reverse=True)
        eatBy1 = [False] * len(reward1)
        for _, _, i in R[:k]:
            eatBy1[i] = True
        res = 0
        for i in range(len(reward1)):
            if eatBy1[i]:
                res += reward1[i]
            else:
                res += reward2[i]
        return res
