from typing import DefaultDict, List, Set, Tuple
from collections import defaultdict


INF = int(1e20)

# 1 <= num <= 10^9
# 1 <= wood.length <= 10^5
# wood[i].length == 2
# 1 <= wood[i][0] <= wood[i][1] <= num
class Solution:
    def buildBridge(self, num: int, wood: List[List[int]]) -> int:
        wood.sort()


print(Solution().buildBridge(num=10, wood=[[1, 2], [4, 7], [8, 9]]))

