"""部门项目分配"""
# N个部门,第i个部门有Ai个雇员,完成一个项目需要K个不同部门的人,
# 问最多能完成多少个项目?
# n<=2e5 Ai<=1e12

# !二分

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))

    def check(mid: int) -> bool:
        """能否凑齐mid个项目 每个部门向尽可能多的项目里去分配人员"""
        res = 0
        for num in nums:
            res += min(mid, num)  # !每个部门可以把人分配到几个项目
        return res >= k * mid

    left, right = 1, n * max(nums) // k
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1

    print(right)
