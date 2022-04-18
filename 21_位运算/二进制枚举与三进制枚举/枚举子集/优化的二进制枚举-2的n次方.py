from typing import List, Set

# 优化的枚举子集
# 也可以dfs(index,visited,curSum)做，也是2^n时间

# 1. dp优化
def getSubsetSums(nums: List[int]) -> List[int]:
    """求每个子集的和 时间复杂度为 O(1+2+4+...+2^(n-1)) = O(2^n)"""
    n = len(nums)
    sums = [0] * (1 << n)
    for i, num in enumerate(nums):
        for pre in range(1 << i):
            cur = sums[pre] + num
            sums[pre | (1 << i)] = cur

    return sums


# 2. 滚动集合更新
# 求所有子集可能的和 O(1+2+4+...+2^(n-1)) = O(2^n)
def getSubsetSum(nums) -> Set[int]:
    res = set([0])
    for num in nums:
        res |= {num + x for x in res} | {num}
    return res


# 3. dfs
def getSubsetSums2(nums: List[int]) -> List[int]:
    """求每个子集的和 时间复杂度为 O(1+2+4+...+2^(n-1)) = O(2^n)"""

    def dfs(index: int, state: int, curSum: int) -> None:
        if index == n:
            res[state] = curSum
            return
        dfs(index + 1, state, curSum)
        dfs(index + 1, state | (1 << index), curSum + nums[index])

    n = len(nums)
    res = [0] * (1 << n)
    dfs(0, 0, 0)
    return res


print(getSubsetSums([1, 2, 3, 1]))
print(getSubsetSum([1, 2, 3, 1]))
print(getSubsetSums2([1, 2, 3, 1]))

