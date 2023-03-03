# n个小朋友 每个小朋友可分0到ai个糖果
# 要分k个糖果 一共有多少种分法 (MOD 1e9+7)
from typing import Tuple
import numpy as np

MOD = int(1e9 + 7)

n, k = map(int, input().split())
A = list(map(int, input().split()))


# 第一个小朋友: (1+x+x^2+...+x^ai) = (1-x^(ai+1))/(1-x)
def gen(i) -> Tuple["np.poly1d", "np.poly1d"]:
    arr1 = [0] * (A[i] + 2)
    arr1[0] = -1
    arr1[-1] = 1
    arr2 = [-1, 1]
    return np.poly1d(arr1), np.poly1d(arr2)


polys = [gen(i) for i in range(n)]
res = np.poly1d([1])
for a, b in polys:
    div, _ = a / b
    # with mod

    res *= div

print(res[k])
