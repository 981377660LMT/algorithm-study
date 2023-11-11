"""游玩游乐园的最大幸福感

给N个数,初始化ans为0,现在有一种操作可以把其中一个数加入ans,
然后该数-1,这样的操作最多进行K次,请问ans的最大值。
n<=1e5,k<=2e9,ai<=2e9

!等价于[1,2,...nums[0],1,2,...nums[1],1,2,...nums[2],...]取k个数的最大和
填平后均摊,二分找到最大的right使得恰好无法凑出k个数
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def arithmeticSum1(first: int, last: int, diff: int) -> int:
    """等差数列求和 first:首项 last:末项 diff:公差"""
    item = (last - first) // diff + 1
    return item * (first + last) // 2


if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))
    # if sum(nums) <= k:
    #     print(sum(arithmeticSum1(1, num, 1) for num in nums))
    #     exit(0)

    def check(mid: int) -> bool:
        """把mid都取完 恰好不能凑齐k个数"""
        res = 0
        for num in nums:
            if num >= mid:
                res += num - mid + 1
        return res < k

    left, right = 1, int(1e10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1

    # !找到最小的left使得left全部取完`恰好不能`凑齐k个数
    res = 0
    for num in nums:
        if num >= right:
            res += arithmeticSum1(num, left, -1)
            k -= num - left + 1
    if k > 0:
        res += k * (left - 1)

    print(res)
