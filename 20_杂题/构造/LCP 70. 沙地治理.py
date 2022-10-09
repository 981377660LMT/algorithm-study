"""
n <= 1000
初始情况下，区域中的所有位置都为沙地，
你需要指定一些子区域种植沙柳树成为绿地，以达到转化整片区域为绿地的最终目的，规则如下：

若两个子区域共用一条边，则视为相邻；
若至少有两片绿地与同一片沙地相邻，则这片沙地也会转化为绿地
转化为绿地的区域会影响其相邻的沙地

构造+找规律 oeis
https://oeis.org/A011848

https://leetcode.cn/problems/XxZZjK/solution/zhao-bu-dao-gui-lu-zhi-jie-ditui-by-pean-cnka/
"""

from copy import deepcopy
from typing import List


MAPPING = {
    0: [],
    1: [[1, 1]],
    2: [[1, 1], [2, 1], [2, 3]],
    3: [[1, 1], [2, 1], [3, 1], [3, 3], [3, 5]],
    4: [[1, 1], [2, 3], [3, 2], [4, 1], [4, 3], [4, 5], [4, 7]],
}


class Solution:
    def sandyLandManagement(self, size: int) -> List[List[int]]:
        def dfs(size: int) -> List[List[int]]:
            if size <= 4:
                return deepcopy(MAPPING[size])

            res = dfs(size - 4)

            # 1
            res.append([size - 3, 1])

            # 2
            for i in range(1, size - 2):
                res.append([size - 2, i * 2 + 1])

            # 3
            res.append([size - 1, 2])

            # 4
            for i in range(size):
                res.append([size, i * 2 + 1])

            return res

        return dfs(size)
