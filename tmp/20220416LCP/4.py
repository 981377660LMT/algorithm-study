from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= time.length == position.length <= 500
# 1 <= time[i] <= 5
# 0 <= position[i] <= 100
class Solution:
    def defendSpaceCity(self, time: List[int], position: List[int]) -> int:
        attack = sorted(zip(time, position))
        print(attack)


print(Solution().defendSpaceCity(time=[1, 2, 1], position=[6, 3, 3]))
