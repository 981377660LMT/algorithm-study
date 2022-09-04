# 顺丰-猜出排列需要的最少次数
# 猜排列游戏

# 小明有一个由1到n的整数组成的排列，
# 他让你来猜出这个排列是什么。
# 你每次可以猜测某一位置的数字，
# 小明会告诉你所猜测的数是“大了”、“小了”或是“正确”。
# !你想知道你在最坏情况下，需要猜测几次，
# 才能在排列的所有位置都得到小明“正确”的回复?
# n<=1e9


# 策略:
# !每次猜中间的数 然后加上 左右两边需要查的次数+每个数需要1次来确定

from functools import lru_cache


@lru_cache(None)
def dfs(remain: int) -> int:
    if remain == 1:
        return 1
    if remain == 2:
        return 3
    if remain & 1:
        return 2 * dfs(remain // 2) + remain
    return dfs(remain // 2) + dfs(remain // 2 - 1) + remain


print(dfs(1))  # 1次
print(dfs(2))  # 3次
print(dfs(5))  # 11次
print(dfs(int(1e9)))
