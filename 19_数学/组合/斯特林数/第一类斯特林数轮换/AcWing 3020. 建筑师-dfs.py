# 1—n个建筑物让他们排列起来，左边与右边分别可以看见A,B个建筑物，问建筑物排列的方案数？
# n<=5e4 A,B<=100
# https://www.acwing.com/solution/content/47769/

# 将最高点n当作分割点 左边看到 A-1个建筑物 右边B-1个建筑物
# 每个小组对应一个圆排列
# 即从n-1个数中选出(A+B-2)个圆排列 再选(A-1)个放在左边
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


def C(n: int, k: int) -> int:
    if n < k:
        return 0
    return (fac(n) * ifac(k) * ifac(n - k)) % MOD


@lru_cache(None)
def cal1(i: int, j: int) -> int:
    """第一类斯特林数:i个人,j个圆排列
    
    i,j<=1000
    - 将该新元素置于一个单独的轮换中 `cal(i - 1, j - 1)
    - 将该元素插入到任何一个现有的轮换中 `(i - 1) * cal(i - 1, j)``
    """
    if i == 0:
        return int(j == 0)
    return (cal1(i - 1, j - 1) + (i - 1) * cal1(i - 1, j)) % MOD


MOD = int(1e9 + 7)


def main(n: int, A: int, B: int) -> int:
    stirling1 = cal1(n - 1, A + B - 2)
    comb = C(A + B - 2, A - 1)
    return stirling1 * comb % MOD


T = int(input())
for _ in range(T):
    n, A, B = map(int, input().split())
    print(main(n, A, B))

