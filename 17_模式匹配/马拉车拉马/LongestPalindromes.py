from typing import List, Tuple


INF = int(1e18)


def longestPalindromesLength(ords: List[int]) -> List[int]:
    """
    对2*n-1个回文中心, 求出每个中心对应的极大回文子串的长度.
    """
    n = len(ords)
    res = [0] * (2 * n - 1)
    palindromes = longestPalindromes(ords)
    for p in palindromes:
        s, e = p
        res[s + e - 1] = e - s
    return res


def longestPalindromes(ords: List[int]) -> List[Tuple[int, int]]:
    """
    给定一个字符串，返回极长回文子串的区间.这样的极长回文子串最多有 2n-1 个.
    """
    n = len(ords)
    m = n * 2 - 1
    sb = [0] * m
    for i in range(n - 1, -1, -1):
        sb[i * 2] = ords[i]
    for i in range(n - 1):
        sb[i * 2 + 1] = INF
    dp = [0] * m
    i, j = 0, 0
    while i < m:
        while i - j >= 0 and i + j < m and sb[i - j] == sb[i + j]:
            j += 1
        dp[i] = j
        k = 1
        while i - k >= 0 and i + k < m and k + dp[i - k] < j:
            dp[i + k] = dp[i - k]
            k += 1
        i += k
        j -= k
    for i in range(m):
        if ((i ^ dp[i]) & 1) == 0:
            dp[i] -= 1
    res = []
    for i in range(m):
        if dp[i] == 0:
            continue
        start = (i - dp[i] + 1) // 2
        end = (i + dp[i] + 1) // 2
        res.append((start, end))
    return res


if __name__ == "__main__":
    # https://judge.yosupo.jp/problem/enumerate_palindromes
    s = input()
    ords = [ord(c) for c in s]
    print(*longestPalindromesLength(ords))
