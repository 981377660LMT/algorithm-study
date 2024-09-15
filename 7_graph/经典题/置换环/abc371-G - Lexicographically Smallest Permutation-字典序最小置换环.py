# G - Lexicographically Smallest Permutation (abc371 G)，字典序最小排列，字典序最小置换环
# https://www.cnblogs.com/Lanly/p/18414807
# 用Python 来算高精度lcm
#
# 由字典序的比较顺序，首先是让第一个数最小
# 那就找第一个数所在的环，遍历一遍，找到数最小的位置，得到一个偏移量 b。记该环的大小为a


from math import gcd
import sys
from typing import List

input = lambda: sys.stdin.readline().rstrip("\r\n")


def collectCycle(nexts: List[int], start: int) -> List[int]:
    """置换环找环.nexts数组中元素各不相同."""
    cycle = []
    cur = start
    while True:
        cycle.append(cur)
        cur = nexts[cur]
        if cur == start:
            break
    return cycle


def minArg(arr: List[int]) -> int:
    minVal, minIdx = arr[0], 0
    for i in range(1, len(arr)):
        if arr[i] < minVal:
            minVal = arr[i]
            minIdx = i
    return minIdx


if __name__ == "__main__":
    N = int(input())
    P = list(map(int, input().split()))  # i->P[i]
    A = list(map(int, input().split()))

    for i in range(N):
        P[i] -= 1

    visited = [False] * N
    period, remain = 1, 0  # !移动次数模period等于remain
    res = [0] * N
    for i in range(N):
        if visited[i]:
            continue
        cycle = collectCycle(P, i)
        for v in cycle:
            visited[v] = True

        # !找到环内最小的数的偏移量
        m = len(cycle)
        count = m // gcd(m, period)  # !可以选择的位置数
        nums = []
        for j in range(count):
            pos = (j * period + remain) % m
            nums.append(A[cycle[pos]])
        offset = minArg(nums)

        remain += offset * period
        period *= count

        for i in range(m):
            res[cycle[i]] = A[cycle[(i + remain) % m]]

    print(*res)
