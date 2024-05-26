import sys
from typing import List


input = lambda: sys.stdin.readline().rstrip("\r\n")

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
    M = 1 << N
    offset = 1
    while offset < M:
        offset <<= 1

    def divide(start: int, end: int) -> List[int]:
        seg = []
        start += offset
        end += offset
        while start < end:
            if start & 1:
                seg.append(start)
                start += 1
            if end & 1:
                end -= 1
                seg.append(end)
            start >>= 1
            end >>= 1
        return seg

    def query(a: int, b: int) -> int:
        print(f"? {a} {b}", flush=True)
        return int(input())

    def output(res: int) -> None:
        print(f"! {res}", flush=True)

    def leftLeaf(x: int) -> int:
        while x < offset:
            x <<= 1
        return x - offset

    def rightLeaf(x: int) -> int:
        while x < offset:
            x = (x << 1) | 1
        return x - offset

    res = []
    for v in divide(L, R + 1):
        left, right = leftLeaf(v), rightLeaf(v)
        i = (right - left + 1).bit_length() - 1
        j = left // 2**i
        res.append(query(i, j))
    output(sum(res) % 100)
