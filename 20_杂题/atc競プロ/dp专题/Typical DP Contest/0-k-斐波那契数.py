"""
N個の0/1の数列のうち1がK個以上連続しない場合の数がk-bonacci numberになることを利用
N=3の場合はhttps://oeis.org/A000073
1番目、N番目は必ず1にならなければいけないので余事象
"""

# !01序列中1不能连续出现k次 -> k-bonacci number
# 特别的，k==2时 斐波那契数列
# k-斐波那契数

# !dp[index][count] 表示前index个数连续1有count个的方案数
# 这里不能存index*count个状态
# 观察转移表达式
# 加1:ndp[j]=dp[j-1] (j>=1)
# 不加1:ndp[0]=dp[0]+dp[1]+...+dp[k-1]

from collections import deque
from typing import List


MOD = int(1e9 + 7)


def kbonacci1(k: int, n: int) -> int:
    """k-bonacci number 第n项

    时间复杂度: O(n)

    >>> a1=a2=...=aK=1
    >>> ai=sum(aj for j in range(i-k,i))


    长为n的01序列中1不能连续出现k次的方案数
    """
    assert k >= 1, n >= 0

    dp = deque([0] * k)
    dp[0] = 1  # TODO1
    sum_ = 1
    for _ in range(n):  # TODO2
        last = dp.pop()
        dp.appendleft(sum_)
        sum_ += sum_ - last
        sum_ %= MOD
    return sum_


"""https://atcoder.jp/contests/tdpc/submissions/15359686"""


def kbonacci2(k: int, n: int) -> int:
    """k-bonacci number 第n项

    时间复杂度: O(k*logk*logn)

    >>> a1=a2=...=aK=1
    >>> ai=sum(aj for j in range(i-k,i))

    长为n的01序列中1不能连续出现k次的方案数
    """
    assert k >= 1, n >= 0

    def linear_recursion_solver(a: List[int], x: List[int], k: int, e0: int, e1: int) -> int:
        def rec(k: int) -> List[int]:
            c = [e0] * m
            if k < m:
                c[k] = e1
                return c[:]
            b = rec(k // 2)
            t = [e0] * (2 * m + 1)
            for i in range(m):
                for j in range(m):
                    t[i + j + (k & 1)] += b[i] * b[j]
                    t[i + j + (k & 1)] %= MOD
            for i in reversed(range(m, 2 * m)):
                for j in range(m):
                    t[i - m + j] += a[j] * t[i]
                    t[i - m + j] %= MOD
            for i in range(m):
                c[i] = t[i]
            return c[:]

        m = len(a)
        c = rec(k)
        res = 0
        for ci, xi in zip(c, x):
            res += ci * xi
            res %= MOD
        return res

    A, C = [1] * k, [1] * k
    return linear_recursion_solver(C[::-1], A, n, 0, 1) % MOD


if __name__ == "__main__":
    print(kbonacci1(2, 2))  # 长度为2的01序列中1不能连续出现2次的方案数 01 10 00 3种
    print(kbonacci2(2, 0))  # 1
    print(kbonacci2(2, 1))  # 1
    print(kbonacci2(2, 2))  # 2
    print(kbonacci2(2, 3))  # 3
    print(kbonacci2(2, 10))  # 3
    # 开头结尾必须为1的情况下 求01序列中1不能连续出现k次的方案数
    # 容斥原理(前后两端无连续k) + k-bonacci number
    n, k = map(int, input().split())
    print((kbonacci1(k, n) - kbonacci1(k, n - 1) * 2 + kbonacci1(k, n - 2)) % MOD)
