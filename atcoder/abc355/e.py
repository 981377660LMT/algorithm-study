import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# この問題は インタラクティブな問題（あなたの作成したプログラムとジャッジシステムが入出力を介して対話を行う形式の問題）です。

# 正整数
# N および
# 0 以上
# 2
# N
#   未満の整数
# L,R(L≤R) が与えられます。 ジャッジシステムは、
# 0 以上
# 99 以下の整数からなる長さ
# 2
# N
#   の数列
# A=(A
# 0
# ​
#  ,A
# 1
# ​
#  ,…,A
# 2
# N
#  −1
# ​
#  ) を隠し持っています。

# あなたの目標は
# A
# L
# ​
#  +A
# L+1
# ​
#  +⋯+A
# R
# ​
#   を
# 100 で割った余りを求めることです。ただし、あなたは数列
# A の要素の値を直接知ることはできません。 その代わりに、ジャッジシステムに対して以下の質問を行うことができます。

# 2
# i
#  (j+1)≤2
# N
#   を満たすように非負整数
# i,j を選ぶ。
# l=2
# i
#  j,r=2
# i
#  (j+1)−1 として
# A
# l
# ​
#  +A
# l+1
# ​
#  +⋯+A
# r
# ​
#   を
# 100 で割った余りを聞く。
# どのような
# A であっても
# A
# L
# ​
#  +A
# L+1
# ​
#  +⋯+A
# R
# ​
#   を
# 100 で割った余りを特定することができる質問回数の最小値を
# m とします。
# m 回以内の質問を行って
# A
# L
# ​
#  +A
# L+1
# ​
#  +⋯+A
# R
# ​
#   を
# 100 で割った余りを求めてください。


if __name__ == "__main__":
    N, L, R = map(int, input().split())

    def divide(start: int, end: int) -> List[Tuple[int, int]]:
        cur = start
        len_ = end - start
        res = []
        for k in range(log, -1, -1):
            if len_ & (1 << k) != 0:
                res.append((k, cur))
                cur += 1 << k
                if cur >= end:
                    return res
        res.append((0, cur))
        return res

    def query(a: int, b: int) -> int:
        print(f"? {a} {b}", flush=True)
        return int(input())

    def output(res: int) -> None:
        print(f"! {res}", flush=True)

    log = (1 << N).bit_length() - 1
    size = N * (log + 1)
    res = []
    print(divide(1, 5 + 1))
    for bit, start in divide(L, R + 1):
        len_ = 1 << bit
        curStart = start
        while start:
            lowbit = start & -start
            start -= lowbit
            curStart += lowbit
    output(sum(res) % 100)
