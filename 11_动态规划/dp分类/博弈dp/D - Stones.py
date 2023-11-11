"""博弈dp

n个石子
两个人轮流取石子,每次可以取nums[i]个石子(nums[i]<=当前剩余石子数)
求先手最多能取多少个石子
n<=1e4 len(nums)<=1e2 nums[0]=1
"""


from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = sorted(set(map(int, input().split())))

    # !贪心是错误的 不一定自己每次取大就是最优解 需要枚举每种情况
    # remain = n
    # res = 0
    # while remain > 0:
    #     pos1 = bisect_right(nums, remain) - 1
    #     res += nums[pos1]
    #     remain -= nums[pos1]
    #     if remain > 0:
    #         pos2 = bisect_right(nums, remain) - 1
    #         remain -= nums[pos2]
    # print(res)

    @lru_cache(None)
    def dfs(remain: int) -> int:
        """先手与后手的最大差值"""
        if remain == 0:
            return 0
        res = -INF
        for num in nums:
            if remain - num >= 0:
                res = max(res, num - dfs(remain - num))
            else:
                break
        return res

    print((n + dfs(n)) // 2)
