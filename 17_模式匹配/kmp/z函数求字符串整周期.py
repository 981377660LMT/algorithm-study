from typing import List


def getZ(string: str) -> List[int]:
    """z算法求字符串公共前后缀的长度

    z[0]=0
    z[i]是s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
    """

    n = len(string)
    Z = [0] * n
    left, right = 0, 0
    for i in range(1, n):
        Z[i] = max(min(Z[i - left], right - i + 1), 0)
        while i + Z[i] < n and string[Z[i]] == string[i + Z[i]]:
            left, right = i, i + Z[i]
            Z[i] += 1
    return Z


def getMinPeriod(s: str) -> int:
    """求字符串的最小周期
    当区间[l+d,r]的哈希值与[l,r-d]的哈希值相等时，那么该区间[l,r]是以 d 为循环节的**
    """
    z = getZ(s)
    for i in range(1, len(s)):
        if len(s) % i == 0 and z[i] == len(s) - i:
            return i
    return -1


if __name__ == "__main__":
    assert getMinPeriod("aabaabaabaab") == 3
    assert getMinPeriod("aaaaa") == 1
    assert getMinPeriod("abcde") == -1
