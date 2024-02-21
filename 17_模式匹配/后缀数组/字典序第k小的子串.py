# https://blog.csdn.net/Elemmir/article/details/50988467
# 所有不相同的子串中字典序第k小的子串
# !二分出排名为K的子串是哪一个后缀的第几个未被计算过的前缀(每个后缀贡献子串数是这个后缀的长度减去其LCP)

from itertools import accumulate
from SuffixArray import SuffixArray


def solve(s: str, k: int) -> str:
    """字典序第k小的子串 k>=1"""
    n = len(s)
    ords = [ord(c) for c in s]
    S = SuffixArray(ords)
    sa, height = S.sa, S.height
    counts = [(n - r - lcp) for r, lcp in zip(sa, height)]
    preSum = [0] + list(accumulate(counts))

    left, right = 0, n
    while left <= right:
        mid = (left + right) // 2
        if preSum[mid] < k:
            left = mid + 1
        else:
            right = mid - 1

    remain = k - preSum[left - 1]  # 排名的k的子串是第left个后缀的第remain个`未出现`过的前缀 需要加上lcp
    remain += height[left - 1]
    start = sa[left - 1]
    end = start + remain
    return s[start:end]


if __name__ == "__main__":
    # 求字符串"EXCITING"的所有不相同子串中字典序排名为20 的子串
    # print(solve("EXCITING", 20))
    # print(solve("EXCITING", 1))
    # print(solve("EXCITING", 35))  # 35个不同子串
    # check with brute force
    from random import randint

    def bruteForce(s: str, k: int) -> str:
        n = len(s)
        subs = set()
        for i in range(n):
            for j in range(i + 1, n + 1):
                subs.add(s[i:j])
        subs = sorted(subs)
        return subs[k - 1]

    while True:
        n = randint(10, 100)
        s = "".join([chr(randint(97, 122)) for _ in range(2 * n)])
        k = randint(1, n * (n + 1) // 2)
        if bruteForce(s, k) != solve(s, k):
            print(s, k)
            break
        else:
            print("ok")
