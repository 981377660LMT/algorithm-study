from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 已知初始 material[i] 表示第 i 种反应物的质量，每次进行实验时，会选出当前 质量最大 的两种反应物进行反应，假设反应物的重量分别为 i 和 j ，且 i <= j。反应的结果如下：

# 如果 i == j，那么两种化学反应物都将被消耗完；
# 如果 i < j，那么质量为 i 的反应物将会完全消耗，而质量为 j 的反应物质量变为 j - i 。
# 最后，最多只会剩下一种反应物，返回此反应物的质量。如果没有反应物剩下，返回 0。


class Solution:
    def lastMaterial(self, material: List[int]) -> int:
        sl = SortedList(material)
        while len(sl) > 1:
            a, b = sl.pop(), sl.pop()
            if a != b:
                sl.add(a - b)
        return sl[0] if sl else 0
