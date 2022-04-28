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
    """
    z = getZ(s)
    return next((i for i in range(len(s)) if z[i] == len(s) - i), -1)


if __name__ == '__main__':
    assert getMinPeriod('aabaabaabaab') == 3
    assert getMinPeriod('aaaaa') == 1
    assert getMinPeriod('abcde') == -1

