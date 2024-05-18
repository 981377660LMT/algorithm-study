# 给你一个数组 nums ，它是 [0, 1, 2, ..., n - 1] 的一个排列。
# 对于任意一个 [0, 1, 2, ..., n - 1] 的排列 perm ，其 分数 定义为：
# score(perm) = |perm[0] - nums[perm[1]]| + |perm[1] - nums[perm[2]]| + ... + |perm[n - 1] - nums[perm[0]]|
# 返回具有 最低 分数的排列 perm 。如果存在多个满足题意且分数相等的排列，则返回其中字典序最小的一个。
#
# !字典序最小：按照字典序从小到大搜索.


from typing import List

INF = int(1e18)


class Solution:
    def findPermutation(self, nums: List[int]) -> List[int]:
        if sorted(nums) == nums:
            return nums

        n = len(nums)
        resCost = INF
        res = [0] * n
        path = [0] * n

        def dfs(index: int, visited: int, pre: int, curRes: int) -> None:
            nonlocal resCost, res
            if curRes >= resCost:
                return
            if index == n:
                curRes += abs(pre - nums[path[0]])
                if curRes < resCost:
                    resCost = curRes
                    res = path[:]
                return
            for i in range(n):
                if visited & (1 << i):
                    continue
                path[index] = i
                dfs(index + 1, visited | (1 << i), i, curRes + abs(pre - nums[i]))

        dfs(1, 1, 0, 0)  # 第一位一定填0
        return res


if __name__ == "__main__":
    # nums = [1,0,2]
    print(Solution().findPermutation([1, 0, 2]))
