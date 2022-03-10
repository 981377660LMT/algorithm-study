# 思路1：对每个数，看他有哪些因子 `n*n^(1/2)`
# 思路2：对每个因子，看他是哪些数的因子 `nlogn`
n = int(input())
factorSums = [0] * (n + 1)
# 1. 因数筛求1-n中每个数的约数有哪些/约数之和
for factor in range(1, n + 1):
    # 约数不能是自己，所以从2开始枚举，如果约数可以是自己，从1开始枚举
    for multi in range(2, n // factor + 1):
        factorSums[factor * multi] += factor
