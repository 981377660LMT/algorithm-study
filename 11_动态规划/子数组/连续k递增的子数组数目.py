# 连续k递增的子数组数目
# 对于每个k（k∈[1，n]），有多少个长度为k的严格递增子数组？

from typing import List

INF = int(1e18)


def countSubarrays(nums: List[int]) -> List[int]:
    """连续k递增的子数组数目(k∈[1,n])"""
    n = len(nums)
    res = [0] * n  # res[i]表示长度为i+1的严格递增子数组个数
    dp, pre = 0, -INF
    for i in range(n):
        if nums[i] > pre:
            dp += 1
        else:
            dp = 1
        res[dp - 1] += 1
        pre = nums[i]
    for i in range(n - 2, -1, -1):
        res[i] += res[i + 1]
    return res


if __name__ == "__main__":
    assert countSubarrays([1, 2, 3, 4]) == [4, 3, 2, 1]
    assert countSubarrays([2, 3, 4, 4, 2]) == [5, 2, 1, 0, 0]
