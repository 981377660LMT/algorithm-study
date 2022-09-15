# 给定一个字符串s
# 求每个字母出现次数为偶数的子序列个数
# C(n,0)+C(n, 2) + C(n, 4) + C(n, 6) + ...  = 2^(n - 1)
# C(n,1)+C(n, 3) + C(n, 5) + C(n, 7) + ...  = 2^(n - 1)

from collections import Counter

MOD = int(1e9 + 7)


def count1(s: str) -> int:
    counter = Counter(s)
    res = 1
    for v in counter.values():
        res *= pow(2, v - 1, MOD)
        res %= MOD
    return res - 1  # 0 is not a subsequence


def count2(s: str) -> int:
    return pow(2, len(s) - len(set(s)), MOD) - 1


if __name__ == "__main__":
    assert count1("ababa") == 7
    assert count2("ababa") == 7
