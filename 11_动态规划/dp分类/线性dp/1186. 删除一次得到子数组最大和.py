# 你可以从原数组中选出一个子数组，并可以决定要不要从中删除一个元素（只能删一次哦），
# （删除后）子数组中至少应当有一个元素，然后该子数组（剩下）的元素总和是所有子数组之中最大的。


from typing import List

# 1 <= arr.length <= 105
# -104 <= arr[i] <= 104


class Solution:
    def maximumSum(self, arr: List[int]) -> int:
        res = -int(1e20)
        remove0, remove1 = -int(1e20), -int(1e20)
        for num in arr:
            # 删0次 删1次
            remove0, remove1 = max(remove0 + num, num), max(remove0, remove1 + num)
            res = max(res, remove0, remove1)
        return res
