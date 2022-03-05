# 给定 n 堆石子，两位玩家轮流操作，
# 每次操作可以从任意一堆石子中拿走任意数量的石子（可以拿完，但不能不拿），
# 最后无法进行操作的人视为失败。

# 问如果两人都采用最优策略，先手是否必胜。

from functools import reduce
from operator import xor

n = int(input())
nums = list(map(int, input().split()))
xor_ = reduce(xor, nums)

# 异或不为0，先手必胜
print('Yes' if xor_ != 0 else 'No')

