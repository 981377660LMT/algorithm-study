from typing import List


# 也可以dfs(index,visited,curSum)做，也是2^n时间
def getSum(nums: List[int]) -> List[int]:
    """求每个子集的和 时间复杂度为 O(1+2+4+...+2^(n-1)) = O(2^n)"""
    n = len(nums)
    sums = [0] * (1 << n)
    for i, num in enumerate(nums):
        for pre in range(1 << i):
            cur = sums[pre] + num
            sums[pre | (1 << i)] = cur

    return sums


print(getSum([1, 2, 3]))

