from itertools import accumulate
from typing import List

# 你最多可以在仓库中放进多少个箱子？
# 1 <= boxes.length, warehouse.length <= 10^5

# 箱子只能从左向右推进仓库中。
# 传递限制
# 总结：从左向右，优先考虑放大的；即仓库从左向右递减


class Solution:
    def maxBoxesInWarehouse(self, boxes: List[int], warehouse: List[int]) -> int:
        warehouse = list(accumulate(warehouse, func=min))
        boxes.sort(reverse=True)  # 最大的放到外面
        boxId, houseId, res = 0, 0, 0
        while boxId < len(boxes) and houseId < len(warehouse):
            if boxes[boxId] <= warehouse[houseId]:
                res += 1
                houseId += 1
            boxId += 1

        return res


print(Solution().maxBoxesInWarehouse(boxes=[4, 3, 4, 1], warehouse=[5, 3, 3, 4, 1]))
# 输出：3
# 我们可以先把高度为 1 的箱子放入 4 号房间，然后再把高度为 3 的箱子放入 1 号、 2 号或 3 号房间，最后再把高度为 4 的箱子放入 0 号房间。
# 我们不可能把所有 4 个箱子全部放进仓库里。
