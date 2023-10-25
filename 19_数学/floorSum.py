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


def modSum(n: int, k: int) -> int:
    """
    余数求和(ModSum/remainderSum)
    ∑k%i (i in [1,n]), 即 sum(k%i for i in [1,n])
    = ∑k-(k/i)*i
    = n*k-∑(k/i)*i
    对于 [l,r] 范围内的 i,k/i 不变，此时 ∑(k/i)*i = (k/i)*∑i = (k/i)*(l+r)*(r-l+1)/2
    """

    def min(a, b):
        return a if a < b else b

    res = n * k
    left, right = 1, 0
    while left <= n:
        h = k // left
        if h > 0:
            right = min(k // h, n)
        else:
            right = n
        w = right - left + 1
        s = (left + right) * w // 2
        res -= h * s
        left = right + 1
    return res


def floorSum2D(n: int, m: int) -> int:
    """
    二维整除分块和。
    ∑{i=1..min(n,m)} floor(n/i)*floor(m/i)
    """

    def min(a, b):
        return a if a < b else b

    res = 0
    left, right = 1, 0
    min_ = min(n, m)
    while left <= min_:
        hn, hm = n // left, m // left
        right = min(n // hn, m // hm)
        w = right - left + 1
        res += hn * hm * w
        left = right + 1
    return res


if __name__ == "__main__":
    print(modSum(*map(int, input().split())))
