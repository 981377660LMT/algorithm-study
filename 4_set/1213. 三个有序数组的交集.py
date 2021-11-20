from typing import List


# 多个数组，多路归并
# 这道题思路跟丑数的思路差不多
class Solution:
    def arraysIntersection(self, arr1: List[int], arr2: List[int], arr3: List[int]) -> List[int]:
        i, j, k = 0, 0, 0
        min_length = min(len(arr1), len(arr2), len(arr3))
        res = []

        while i < min_length and j < min_length and k < min_length:
            # 三个数值相同，指针都动一位
            if arr1[i] == arr2[j] == arr3[k]:
                res.append(arr1[i])
                i, j, k = i + 1, j + 1, k + 1
                continue

            # 不全相同，则最小数的指针移动一位
            min_val = min(arr1[i], arr2[j], arr3[k])
            if arr1[i] == min_val:
                i += 1
            if arr2[j] == min_val:
                j += 1
            if arr3[k] == min_val:
                k += 1

        return res


print(Solution().arraysIntersection([1, 2, 3, 4, 5], [1, 2, 5, 7, 9], [1, 3, 4, 5, 8]))
