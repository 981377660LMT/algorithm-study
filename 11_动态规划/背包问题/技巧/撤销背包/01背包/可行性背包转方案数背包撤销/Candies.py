# https://www.luogu.com.cn/problem/P6808
# 你需要修改序列中的一个数 P 为 Q，使得尽可能多的整数能够被表示出来。

from typing import List, Tuple
from Knapsack01Removable import Knapsack01Removable


MOD = int(1e9 + 7)  # 大素数


def candies(nums: List[int]) -> List[Tuple[int, int]]:
    res = 0
    sum_ = sum(nums)
    dp = Knapsack01Removable(sum_, MOD)
    for v in nums:
        dp.add(v)
    for v in nums:
        dp.remove(v)

        dp.add(v)
    return res


if __name__ == "__main__":
    n = map(int, input().split())
    nums = list(map(int, input().split()))
    p, q = candies(nums)
    print(p, q)
