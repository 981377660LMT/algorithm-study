"""
交互

!1. 给定st表的n,输出需要预处理的区间
    长为1的区间:[0,0],[1,1],...,[n-1,n-1]
    长为2的区间:[0,1],[1,2],...,[n-2,n-1]
    长为4的区间:[0,3],[1,4],...,[n-4,n-1]
    ...
    dp[i][j]表示区间[j,j+2**i-1]的贡献
    区间个数上界为n*bit_length(n)

!2. 给定每次的查询区间[L,R],输出分解成的两个区间
    k = (R-L+1).bit_length()-1
    [L,L+(1<<k)-1] 与 [R-(1<<k)+1,R] 对应 st表的dp[k][L]与dp[k][R-(1<<k)+1]
"""


def unionOfTwoSets() -> None:
    n = int(input())

    # phase1
    size = n.bit_length()
    dpIntervals = dict()
    for i in range(n):  # init
        dpIntervals[(i, i)] = len(dpIntervals)
    for i in range(1, size):  # fill
        for j in range(n - (1 << i) + 1):
            dpIntervals[(j, j + (1 << i) - 1)] = len(dpIntervals)
    print(len(dpIntervals), flush=True)
    for left, right in dpIntervals:
        print(left + 1, right + 1, flush=True)

    # phase2
    q = int(input())
    for _ in range(q):
        left, right = map(int, input().split())
        left, right = left - 1, right - 1
        k = (right - left + 1).bit_length() - 1
        interval1 = (left, left + (1 << k) - 1)
        interval2 = (right - (1 << k) + 1, right)
        print(dpIntervals[interval1] + 1, dpIntervals[interval2] + 1, flush=True)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    unionOfTwoSets()
