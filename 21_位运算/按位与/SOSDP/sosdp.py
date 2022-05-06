from functools import lru_cache
from math import ceil, log2


@lru_cache(None)
def sosdp(state: int) -> int:
    """state真子集的贡献"""
    if state == 0:
        return 1
    res = 0
    for i in range(20):
        if (state >> i) & 1:
            res += sosdp(state ^ (1 << i))
    return res


print(log2(1e6))  # 取20位即可
print(sosdp(5))

####################################
n = ceil(log2(1e5))
upper = 1 << n
preSum = [0] * upper
for i in range(upper):
    preSum[i] = bin(i).count('1')
for i in range(n):
    for state in range(upper):
        if (state >> i) & 1:
            preSum[state] += preSum[state ^ (1 << i)]
print(preSum[5])  # 101 100 001 000 四个子集的贡献
