# 珠宝的价值分别为 2,3,…,n+1。
# 给这些珠宝染色，使得一件珠宝的价格是另一件珠宝的价格的`质因子`时，两件珠宝的颜色不同。
# 并且，华生要求他使用的颜色数尽可能少。
# 请帮助夏洛克完成这个简单的任务。

# 输出格式
# 第一行一个整数 k，表示所使用的颜色数；
# 第二行 n 个整数，表示第 1 到第 n 件珠宝被染成的颜色。

'''
只要序列里面有合数，一定需要两种颜色，不论数有多少，至多需要两种颜色
特殊情况，当序列只有质数时候，只需要一种颜色
'''

from typing import List


def getPrimes(upper: int) -> List[int]:
    """筛选出1-upper中的质数"""
    visited = [False] * (upper + 1)
    res = []
    for num in range(2, upper + 1):
        if visited[num]:
            continue
        res.append(num)
        for multi in range(num * num, upper + 1, num):
            visited[multi] = True

    return res


n = int(input())
print(2 if n + 1 > 3 else 1)
primes = set(getPrimes(n + 10))

for i in range(2, n + 2):
    if i in primes:
        print(1, end=' ')
    else:
        print(2, end=' ')


# 一定是一个二分图 质数、合数
# 有边就是2
