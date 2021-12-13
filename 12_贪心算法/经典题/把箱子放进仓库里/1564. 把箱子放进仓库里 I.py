from typing import List

# 你最多可以在仓库中放进多少个箱子？
# 1 <= boxes.length, warehouse.length <= 10^5

# 箱子只能从左向右推进仓库中。

# 总结：从左向右，优先考虑放大的；即仓库从左向右递减
class Solution:
    def maxBoxesInWarehouse(self, boxes: List[int], warehouse: List[int]) -> int:
        boxes = sorted(boxes, reverse=True)
        res, box_id = 0, 0

        for w in warehouse:
            while box_id < len(boxes) and w < boxes[box_id]:
                box_id += 1
            if box_id == len(boxes):
                break
            box_id += 1
            res += 1

        return res


print(Solution().maxBoxesInWarehouse(boxes=[4, 3, 4, 1], warehouse=[5, 3, 3, 4, 1]))
# 输出：3
# 我们可以先把高度为 1 的箱子放入 4 号房间，然后再把高度为 3 的箱子放入 1 号、 2 号或 3 号房间，最后再把高度为 4 的箱子放入 0 号房间。
# 我们不可能把所有 4 个箱子全部放进仓库里。
