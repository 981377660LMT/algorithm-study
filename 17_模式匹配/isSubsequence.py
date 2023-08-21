# O(n+m)子序列匹配


from typing import Any, Sequence


def isSubsequence(longer: str, shorter: str) -> bool:
    if len(shorter) > len(longer):
        return False
    it = iter(longer)
    return all(need in it for need in shorter)


def isSubsequence2(longer: Sequence[Any], shorter: Sequence[Any]) -> bool:
    """判断shorter是否是longer的子序列"""
    if len(shorter) > len(longer):
        return False
    if len(shorter) == 0:
        return True
    i, j = 0, 0
    while i < len(longer) and j < len(shorter):
        if longer[i] == shorter[j]:
            j += 1
            if j == len(shorter):
                return True
        i += 1
    return False


if __name__ == "__main__":
    assert isSubsequence("aabbccdd", "abc")
    assert isSubsequence2("aabbccdd", "abc")

    # 2825. 循环增长使字符串子序列等于另一个字符串
    # https://leetcode.cn/problems/make-string-a-subsequence-using-cyclic-increments/description/
    # 给你一个下标从 0 开始的字符串 str1 和 str2 。
    # 一次操作中，你选择 str1 中的若干下标。对于选中的每一个下标 i ，你将 str1[i] 循环 递增，变成下一个字符。也就是说 'a' 变成 'b' ，'b' 变成 'c' ，以此类推，'z' 变成 'a' 。
    # 如果执行以上操作 至多一次 ，可以让 str2 成为 str1 的子序列，请你返回 true ，否则返回 false 。
    # 注意：一个字符串的子序列指的是从原字符串中删除一些（可以一个字符也不删）字符后，剩下字符按照原本先后顺序组成的新字符串。

    class Solution:
        def canMakeSubsequence(self, str1: str, str2: str) -> bool:
            if len(str1) < len(str2):
                return False
            cand = []
            for c in str1:
                cand.append(c)
                next_ = chr(ord(c) + 1) if c != "z" else "a"
                cand.append(next_)
            cand = "".join(cand)
            return isSubsequence(cand, str2)
