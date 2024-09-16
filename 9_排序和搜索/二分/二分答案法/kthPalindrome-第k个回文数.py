# kthPalindrome-第k个回文数(0也是回文数)
#
# 1 桁の回文数：
# 1,2,3,…,9 の 9 個
# 2 桁の回文数：
# 11,22,33,…,99 の 9 個
# 3 桁の回文数：
# 101,111,121,…,999 の 90 個
# 4 桁の回文数：
# 1001,1111,1221,…,9999 の 90 個


def kthPalindrome(k: int) -> str:
    """返回第k(k>=1)个回文数(0也是回文数)."""
    if k == 1:
        return "0"

    d = 0
    count = 1
    preCount = 0
    while count < k:
        preCount = count
        count += 9 * 10 ** (d // 2)
        d += 1

    half = str(10 ** ((d - 1) // 2) + (k - preCount - 1))
    if d % 2 == 0:
        return half + half[::-1]
    else:
        return half + half[:-1][::-1]


# D - Palindromic Number
# https://atcoder.jp/contests/abc363/tasks/abc363_d
if __name__ == "__main__":
    K = int(input())
    print(kthPalindrome(K))
