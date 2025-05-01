# https://leetcode.cn/problems/count-good-triplets/solutions/3622921/liang-chong-fang-fa-bao-li-mei-ju-qian-z-apcv/
# 类似前缀和的思想，这等价于：
# right 中的 ≤x+c 的元素个数，减去 right 中的 <x−c 的元素个数。


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


if __name__ == "__main__":

    def bf(arr1, arr2, k):
        res = 0
        for i in arr1:
            for j in arr2:
                if abs(i - j) <= k:
                    res += 1
        return res

    def bf2(arr1, arr2, k):
        import bisect

        res = 0
        for v in arr1:
            left = bisect.bisect_left(arr2, v - k)
            right = bisect.bisect_right(arr2, v + k)
            res += right - left
        return res

    import random

    for _ in range(100):
        arr1 = sorted([random.randint(0, 100) for _ in range(100)])
        arr2 = sorted([random.randint(0, 100) for _ in range(100)])
        k = random.randint(0, 100)
        assert countPairsWithinDifference(arr1, arr2, k) == bf(arr1, arr2, k) == bf2(arr1, arr2, k)
    print("passed")
