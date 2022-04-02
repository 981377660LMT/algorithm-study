# 分数为1,2,...,n
# 还有最后一轮比赛，求可能赢的人数

# 贪心大多数需要排序
class Solution:
    def solve(self, nums):
        if not nums:
            return 0
        nums = sorted(nums)

        # 贪心的给第一名1分，第二2分，并维护max
        # 最后检查是否每个数有机会超过max
        n = len(nums)
        best = 0
        for i in range(n):
            best = max(best, nums[i] + n - i)

        res = 0
        for i in range(n):
            if nums[i] + n >= best:
                res += 1
        return res


print(Solution().solve(nums=[8, 7, 10, 11]))
