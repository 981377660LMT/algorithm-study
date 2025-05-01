# 1534. 统计好三元组
# https://leetcode.cn/problems/count-good-triplets/description/
# 如果三元组 (arr[i], arr[j], arr[k]) 满足下列全部条件，则认为它是一个 好三元组 。
#
# 0 <= i < j < k < arr.length
# |arr[i] - arr[j]| <= a
# |arr[j] - arr[k]| <= b
# |arr[i] - arr[k]| <= c
# 其中 |x| 表示 x 的绝对值。
#
# 返回 好三元组的数量 。
#
# 排序 + 枚举中间 + 三指针 O(n^2)


from typing import List


def countPairsWithinDifference(arr1: List[int], arr2: List[int], k: int) -> int:
    """给定两个有序数组，从两个数组中各选一个数，计算绝对差 ≤k 的数对个数."""
    res = 0
    left, right = 0, 0
    for v in arr1:
        while right < len(arr2) and arr2[right] <= v + k:
            right += 1
        while left < len(arr2) and arr2[left] < v - k:
            left += 1
        res += right - left
    return res


class Solution:
    def countGoodTriplets(self, arr: List[int], a: int, b: int, c: int) -> int:
        order = sorted(range(len(arr)), key=lambda i: arr[i])
        res = 0
        for j in order:
            y = arr[j]
            # 根据order遍历，left、right 天然有序
            left = [arr[i] for i in order if i < j and abs(arr[i] - y) <= a]
            right = [arr[k] for k in order if k > j and abs(arr[k] - y) <= b]
            res += countPairsWithinDifference(left, right, c)
        return res
