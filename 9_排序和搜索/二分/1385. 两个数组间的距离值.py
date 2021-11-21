# 对于元素 arr1[i] ，恒有 |arr1[i]-arr2[j]| > d 。
# 计算arr1中这样的元素个数

from typing import List
import bisect


class Solution:
    def __findTheDistanceValue(self, arr1: List[int], arr2: List[int], d: int) -> int:
        return sum(all(abs(a1 - a2) > d for a2 in arr2) for a1 in arr1)

    def findTheDistanceValue(self, arr1: List[int], arr2: List[int], d: int) -> int:
        arr2.sort()
        res = 0

        for num in arr1:
            index = bisect.bisect_left(arr2, num)
            # 比较左右
            is_right_valid = index == len(arr2) or arr2[index] - num > d
            is_left_valid = index == 0 or num - arr2[index - 1] > d
            if is_left_valid and is_right_valid:
                res += 1

        return res


print(Solution().findTheDistanceValue([4, 5, 8], [10, 9, 1, 8], 2))
