from typing import List

# 箱子可以从任意方向（左边或右边）推入仓库中。
# 总结：里面的要塞小的箱子；否则外面塞小的箱子会把高的位置他挡住
# 即把仓库变成从左右向中间递减的数组


class Solution:
    def maxBoxesInWarehouse(self, boxes: List[int], warehouse: List[int]) -> int:
        arr = [0x7FFFFFFF] * len(warehouse)
        min_val = 0x7FFFFFFF

        for i, val in enumerate(warehouse):
            min_val = min(min_val, val)
            arr[i] = min_val

        min_val = 0x7FFFFFFF
        for i in range(len(warehouse) - 1, -1, -1):
            min_val = min(min_val, warehouse[i])
            arr[i] = max(arr[i], min_val)

        arr.sort(reverse=True)
        boxes.sort(reverse=True)

        res = 0
        i = 0
        for box in boxes:
            if i == len(arr):
                break

            if box <= arr[i]:
                res += 1
                i += 1

        return res


print(Solution().maxBoxesInWarehouse(boxes=[1, 2, 2, 3, 4], warehouse=[3, 4, 1, 2]))
