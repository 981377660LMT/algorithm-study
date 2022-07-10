# 给定整数 N，试把阶乘 N! 分解质因数，按照算术基本定理的形式输出分解结果中的 pi 和 ci 即可。
# 3≤N≤1e6
"""
先把1-1000000中所有质数全部筛出来
然后统计每一个质数出现的次数，假设
当前统计的质数是P

P在N阶乘中出现的次数就是看1-N中所有数字
出现P因子的次数总和
根据勒让德定理就是 N//P + N//(P^2) + N//(p^3) + ......


"""
from collections import Counter
from typing import List


# def getPrimeFactors(n: int) -> Counter:
#     """返回 n 的所有质数因子"""
#     res = Counter()
#     upper = floor(n ** 0.5) + 1
#     for i in range(2, upper):
#         while n % i == 0:
#             res[i] += 1
#             n //= i

#     # 注意考虑本身
#     if n > 1:
#         res[n] += 1
#     return res


def getPrimes(upper: int) -> List[int]:
    """筛选出1-upper中的质数"""
    visited = [False] * (upper + 1)
    for num in range(2, upper + 1):
        if visited[num]:
            continue
        for multi in range(num * num, upper + 1, num):
            visited[multi] = True

    return [num for num in range(2, upper + 1) if not visited[num]]


primes = getPrimes(int(1e6))

n = int(input())
counter = Counter()
for p in primes:
    multi = 1
    while p**multi <= n:
        counter[p] += n // (p**multi)
        multi += 1

res = sorted((key, count) for key, count in counter.items())
for key, count in res:
    print(key, count)
