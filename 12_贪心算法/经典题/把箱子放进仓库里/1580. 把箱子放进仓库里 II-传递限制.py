from itertools import accumulate
from typing import List

# 箱子可以从任意方向（左边或右边）推入仓库中。

# 传递限制
# 即把仓库变成从左右向中间递减的数组
# 从仓库最低处分为两段，从而转化为1564. 把箱子放进仓库里 I


class Solution:
    def maxBoxesInWarehouse(self, boxes: List[int], warehouse: List[int]) -> int:
        preMin = list(accumulate(warehouse, func=min))
        sufMin = list(accumulate(warehouse[::-1], func=min))[::-1]
        warehouse = list(max(a, b) for a, b in zip(preMin, sufMin))  # 变为V形

        warehouse.sort(reverse=True)
        boxes.sort(reverse=True)

        boxId, houseId, res = 0, 0, 0
        while boxId < len(boxes) and houseId < len(warehouse):
            if boxes[boxId] <= warehouse[houseId]:
                res += 1
                houseId += 1
            boxId += 1

        return res


print(Solution().maxBoxesInWarehouse(boxes=[1, 2, 2, 3, 4], warehouse=[3, 4, 1, 2]))
