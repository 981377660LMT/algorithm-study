# z函数-扩展kmp

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


def zAlgo(string: str) -> List[int]:
    """z算法求字符串每个后缀与原串的最长公共前缀长度

    z[0]=0
    z[i]是s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
    """

    n = len(string)
    z = [0] * n
    left, right = 0, 0
    for i in range(1, n):
        z[i] = max2(min2(z[i - left], right - i + 1), 0)
        while i + z[i] < n and string[z[i]] == string[i + z[i]]:
            left, right = i, i + z[i]
            z[i] += 1
    return z


if __name__ == "__main__":

    class Solution1:
        # 求的就是扩展 KMP（Z 数组）的所有元素之和
        def sumScores(self, s: str) -> int:
            z = zAlgo(s)
            return sum(z) + len(s)

    # 100203. 将单词恢复初始状态所需的最短时间 II
    # https://leetcode.cn/problems/minimum-time-to-revert-word-to-initial-state-ii/description/
    # 给你一个下标从 0 开始的字符串 word 和一个整数 k 。
    # 在每一秒，你必须执行以下操作：
    # 移除 word 的前 k 个字符。
    # 在 word 的末尾添加 k 个任意字符。
    # 注意 添加的字符不必和移除的字符相同。但是，必须在每一秒钟都执行 两种 操作。
    # 返回将 word 恢复到其 初始 状态所需的 最短 时间（该时间必须大于零）。
    #
    # !如果只操作一次，就能让 s恢复成其初始值，意味着什么？
    # !由于可以往 s 的末尾添加任意字符，这意味着只要 s[k:] 是 s 的前缀，就能让 s 恢复成其初始值
    class Solution:
        def minimumTimeToInitialState(self, word: str, k: int) -> int:
            Z = zAlgo(word)
            n = len(word)
            res = 1
            ptr = k
            while ptr < n:
                remain = n - ptr
                if Z[ptr] >= remain:
                    return res
                ptr += k
                res += 1
            return res

    print(Solution().minimumTimeToInitialState(word="abacaba", k=3))
