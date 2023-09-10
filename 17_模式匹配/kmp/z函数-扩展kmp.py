from typing import List


def zAlgo(string: str) -> List[int]:
    """z算法求字符串每个后缀与原串的最长公共前缀长度

    z[0]=0
    z[i]是s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
    """

    n = len(string)
    z = [0] * n
    left, right = 0, 0
    for i in range(1, n):
        z[i] = max(min(z[i - left], right - i + 1), 0)
        while i + z[i] < n and string[z[i]] == string[i + z[i]]:
            left, right = i, i + z[i]
            z[i] += 1
    return z


if __name__ == "__main__":

    class Solution:
        # 求的就是扩展 KMP（Z 数组）的所有元素之和
        def sumScores(self, s: str) -> int:
            z = zAlgo(s)
            return sum(z) + len(s)
