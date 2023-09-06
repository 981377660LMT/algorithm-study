from typing import List, Set


def subsetSum(nums: List[int]) -> List[int]:
    """O(2^n)求所有子集对应的和"""
    n = len(nums)
    res = [0] * (1 << n)
    for i in range(n):
        for pre in range(1 << i):
            res[pre + (1 << i)] = res[pre] + nums[i]
    return res


def subsetSum2(nums: List[int]) -> Set[int]:
    """O(2^n)求所有非空子集的可能和"""
    dp = set()
    for cur in nums:
        dp |= {(cur + pre) for pre in (dp | {0})}
    return dp


if __name__ == "__main__":
    print(subsetSum([1, 2, 3]))
    print(subsetSum2([1, 2, 4]))
