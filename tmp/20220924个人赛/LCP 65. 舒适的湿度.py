# n<=1000 nums[i]<=1000
# !还是需要看数据量猜dp方法
# !dp[i][val]

from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= operate.length <= 1000
# 1 <= operate[i] <= 1000
class Solution:
    def a(self, s: str) -> List[str]:
        ...
