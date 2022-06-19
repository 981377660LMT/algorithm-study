# 请你返回可在 `一次交换`（交换两数字 A[i] 和 A[j] 的位置）后得到的、
# 按字典序排列小于 A 的最大可能排列。

# !注意这里只交换一次 所以不是prePermutation

# !从后往前找，找到第一个arr[i]>arr[i+1]的位置记为i，然后从这个位置往后找，
# !找比arr[i]小的最大元素位置记为j，交换两者即可
from typing import List


class Solution:
    def prevPermOpt1(self, arr: List[int]) -> List[int]:
        n = len(arr)
        left = next((i for i in range(n - 2, -1, -1) if arr[i] > arr[i + 1]), -1)
        if left == -1:
            return arr

        right = max(
            (i for i in range(left, n) if arr[i] < arr[left]), key=lambda x: arr[x], default=-1
        )
        if right == -1:
            return arr

        arr[left], arr[right] = arr[right], arr[left]
        return arr

