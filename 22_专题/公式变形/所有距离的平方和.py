# https://atcoder.jp/contests/abc194/tasks/abc194_c
# 求所有组合的平方和
# n<=3e5

# 所有组合平方和等于product平方和的一半

# 1/2*∑∑(ai-aj)^2
# =1/2*n*∑(ai^2-2∑∑aiaj+aj^2)
# =n*∑ai^2 - (∑ai)^2


from typing import List


def squareSum(nums: List[int]):
    """所有组合的平方和"""
    n = len(nums)
    sum_ = sum(nums)
    return n * sum(num * num for num in nums) - sum_ * sum_


print(squareSum([1, 2, 3]))

# 縦を横に見る（シグマの内側外側を逆転させる）
# 倍数を約数から見る
# 事前計算
# 最終結果だけに着目する
# 操作に決まった順番がある
