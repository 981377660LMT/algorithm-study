# 求∑floor((A*i+B)/M) 的值 (i=0-N-1)
# N,M<=1e9

# !floor sum
# floor_sumは格子点の数え上げ問題に帰着出来る 数格子问题
# 愚直にやるとTLE


def floorSum(N: int, M: int, A: int, B: int) -> int:
    """求∑floor((A*i+B)/M) 的值 (i=0-N-1)"""
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


T = int(input())
for _ in range(T):
    N, M, A, B = map(int, input().split())
    print(floorSum(N, M, A, B))

