from typing import List


def distinct_subsets(nums: List[int]):
    """枚举所有子集，数组中可能有重复元素，结果中不能有重复的子集。"""
    nums = sorted(nums)
    subset = []

    def dfs(start: int):
        yield subset[:]
        for i in range(start, len(nums)):
            if i > start and nums[i] == nums[i - 1]:
                continue
            subset.append(nums[i])
            yield from dfs(i + 1)
            subset.pop()

    yield from dfs(0)


if __name__ == "__main__":

    class Solution:
        def subsetsWithDup(self, nums: List[int]) -> List[List[int]]:
            return list(distinct_subsets(nums))
