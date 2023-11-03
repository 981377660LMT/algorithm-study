# https://atcoder.jp/contests/abc056/tasks/arc070_b
# 给定一个大小为 n 的正整数集合，一个元素和不小于 k 的非空子集被称为优秀的，
# 问有多少个元素满足把它改成 0 而不会影响优秀子集的个数。
# n,k<=5000
#
# 解法：
# 1. >=k的数不是可有可无的。
# 2. 对每个物品撤销，如果撤销后[k-nums[i],k-1]内所有dp值为0，那这个数就是可有可无的数(新增不会改变优秀个数)。
# 3. 分析一下体积上限:2*k

from typing import List
from Knapsack01Removable import Knapsack01Removable


def noNeed(nums: List[int], k: int) -> int:
    MOD = int(1e9 + 7)  # 大素数
    res = 0
    dp = Knapsack01Removable(2 * k, MOD)
    for v in nums:
        dp.add(v)
    for v in nums:
        if v >= k:
            continue
        dp.remove(v)

        ok = True
        for i in range(k - v, k):
            if dp.query(i) > 0:
                ok = False
                break
        if ok:
            res += 1
        dp.add(v)
    return res


if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))
    print(noNeed(nums, k))
