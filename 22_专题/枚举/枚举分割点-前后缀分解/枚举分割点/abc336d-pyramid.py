# https://atcoder.jp/contests/abc336/tasks/abc336_d
# 操作成为金字塔数组的最大长度k
# 操作1: 某个数-1
# 操作2: 删除首或尾
# !这里金字塔数组的定义是: 1,2,..k-1,k,k-1,..2,1
# nums[i]>=1

from typing import List


def pyramid(nums: List[int]) -> int:
    def makeDp(seq: List[int]) -> List[int]:
        n = len(seq)
        dp = [1] * n
        for i in range(1, n):
            dp[i] = min(dp[i - 1] + 1, seq[i])
        return dp

    # pre[i] 表示以nums[i]为右端点的左侧的最长金字塔数组长度
    # suf[i] 表示以nums[i]为左端点的右侧的最长金字塔数组长度
    pre, suf = makeDp(nums), makeDp(nums[::-1])[::-1]
    res = 0
    for i in range(len(nums)):
        res = max(res, min(pre[i], suf[i]))
    return res


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(pyramid(nums))
