# 验证所有小于一百万的偶数能否满足哥德巴赫猜想。
# 任意一个大于 4 的偶数都可以拆成两个奇素数之和。
'''
先用线性筛法将1000000以内的质数全部找出来
然后对于每一个n, 枚举比其小的质数即可
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


primes = set(getPrimes(1000000))

while True:
    n = int(input())
    if n == 0:
        break
    for p1 in primes:
        if (n - p1) in primes:
            print(f'{n} = {p1} + {n-p1}')
            break

