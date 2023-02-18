# https://judge.yosupo.jp/problem/sum_of_floor_of_linear
# 求∑floor((A*i+B)/M) 的值 (i=0-N-1)
# N,M<=1e9
# O(logn+m+a+b)

# !floor sum
# floor_sumは格子点の数え上げ問題に帰着出来る 数格子问题
# 愚直にやるとTLE


def floorSum(N: int, M: int, A: int, B: int) -> int:
    """
    ```
    sum((A*i+B)//M for i in range(N)
    ```
    """
    res = 0
    while True:
        if A >= M or A < 0:
            res += N * (N - 1) * (A // M) // 2
            A %= M
        if B >= M or B < 0:
            res += N * (B // M)
            B %= M

        yMax = A * N + B
        if yMax < M:
            break
        N, B, M, A = yMax // M, yMax % M, A, M
    return res


# verify www.codechef.com/viewsolution/36222026
# count x : ax + b mod m 0<= yr<=m, 0 <= x < xr
def range_count(a: int, b: int, m: int, xr: int, yr: int) -> int:
    return floorSum(xr, m, a, b + m) - floorSum(xr, m, a, b + m - yr)


T = int(input())
for _ in range(T):
    N, M, A, B = map(int, input().split())
    print(floorSum(N, M, A, B))
