# quotient_range/EnumerateQuotients/enumerateFloor
# 数论分块-商列举


from typing import List, Tuple


def floorRange(n: int) -> List[Tuple[int, int, int]]:
    """
    将 [1,n] 内的数分成O(2*sqrt(n))段, 每段内的 n//i 相同

    Args:
        n (int): n>=1

    Returns:
        List[Tuple[int,int,int]]:
        每个元素为(left,right,div)
        表示 left <= i <= right 内的 n//i == div
    """
    if n <= 0:
        return []
    res = []
    m = 1
    while m * m <= n:
        res.append((m, m, n // m))
        m += 1
    for i in range(m, 0, -1):
        left = n // (i + 1) + 1
        right = n // i
        if left <= right and res and res[-1][1] < left:
            res.append((left, right, n // left))
    return res


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc414/tasks/abc414_e
    # 给定正整数 N，求满足条件的三元组 (a, b, c) 个数 mod 998244353，其中 1≤a,b,c≤N 且三者互不相同，且 a mod b = c。
    def abc414_e() -> None:
        MOD = 998244353
        N = int(input())

        # sum1 = sum_{i=1..N} floor(N/i) mod MOD
        sum1 = 0
        for l, r, d in floorRange(N):
            count = r - l + 1
            sum1 += (count * d) % MOD
            sum1 %= MOD

        total = N * (N + 1) // 2 % MOD

        res = (total - sum1) % MOD
        print(res)

    def yukicoder1573() -> None:
        # https://yukicoder.me/problems/no/1573
        # 约数总和
        MOD = 998244353
        n, m = map(int, input().split())
        res = 0
        for left, right, div in floorRange(n):
            right += 1
            lower = left
            higher = min(right - 1, m)
            if lower > higher:
                break
            x = div * (div + 1) // 2 + div
            y = (lower + higher) * (higher - lower + 1) // 2
            res += x * y
            res %= MOD
        print(res)

        # https://judge.yosupo.jp/problem/enumerate_quotients
        n = int(input())
        res = floorRange(n)
        print(len(res))
        for left, right, div in res[::-1]:
            print(div, end=" ")

    abc414_e()
