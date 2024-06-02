from collections import defaultdict
from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 長さ
# N の数列
# A=(A
# 1
# ​
#  ,…,A
# N
# ​
#  ) が与えられます。

# i=1
# ∑
# N−1
# ​

# j=i+1
# ∑
# N
# ​
#  ⌊
# min(A
# i
# ​
#  ,A
# j
# ​
#  )
# max(A
# i
# ​
#  ,A
# j
# ​
#  )
# ​
#  ⌋ を求めてください。

# ただし、
# ⌊x⌋ は
# x 以下の最大の整数を表します。例えば、
# ⌊3.14⌋=3、
# ⌊2⌋=2 です。


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    A.sort()

    diff = [0] * (max(A) + 1)
    counter = defaultdict(int)
    for a in A:
        counter[a] += 1

    for divisor in counter:
        for dividend in range(divisor, len(diff), divisor):
            diff[dividend] += counter[divisor]
    contribution = list(accumulate(diff))
    res = sum(contribution[num] for num in A) - sum(c * (c + 1) // 2 for c in counter.values())
    print(res)
