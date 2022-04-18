from itertools import chain, combinations
from typing import Any, Collection, Generator, List


# 1. dfs
class Solution1:
    def subsets(self, nums: List[int]) -> List[List[int]]:
        def dfs(index: int, path: List[int]):
            if index == len(nums):
                yield path[:]
                return
            path.append(nums[index])
            yield from dfs(index + 1, path)
            path.pop()
            yield from dfs(index + 1, path)

        return list(dfs(0, []))


# 2. powerset 顺序枚举
def powerset(collection: Collection[Any], isProperSubset=False):
    """求(真)子集,时间复杂度O(n*2^n)

    默认求所有子集
    """
    upper = len(collection) if isProperSubset else len(collection) + 1
    return chain.from_iterable(combinations(collection, n) for n in range(upper))


# 3. 枚举+check 顺序枚举
class Solution2:
    def subsets(self, nums: List[int]) -> List[List[int]]:
        def gen() -> Generator[List[int], None, None]:
            n = len(nums)
            for state in range(1 << n):
                cur = []
                for i in range(n):
                    if state & (1 << i):
                        cur.append(nums[i])
                yield cur

        return list(gen())


# 4. 滚动更新
class Solution3:
    def subsets(self, nums: List[int]) -> List[List[int]]:
        res = [[]]
        for num in nums:
            cur = []
            for pre in res:
                cur.append(pre + [num])
            res += cur
        return res
