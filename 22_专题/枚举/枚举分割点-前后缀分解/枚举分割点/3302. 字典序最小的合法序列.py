# 3302. 字典序最小的合法序列
# https://leetcode.cn/problems/find-the-lexicographically-smallest-valid-sequence/description/
#
# 给你两个字符串 word1 和 word2 。
# 如果一个字符串 x 修改 至多 一个字符会变成 y ，那么我们称它与 y 几乎相等 。
# 如果一个下标序列 seq 满足以下条件，我们称它是 合法的 ：
# 下标序列是 升序 的。
# 将 word1 中这些下标对应的字符 按顺序 连接，得到一个与 word2 几乎相等 的字符串。
# !请你返回一个长度为 word2.length 的数组，表示一个字典序最小的 合法 下标序列。如果不存在这样的序列，请你返回一个 空 数组。
# 注意 ，答案数组必须是字典序最小的下标数组，而 不是 由这些下标连接形成的字符串。
#
# !最多失配一次(失配一个字符)的子序列匹配


from typing import Any, List, Sequence


def matchSubsequence(longer: Sequence[Any], shorter: Sequence[Any]) -> List[int]:
    """返回 longer 的每个前缀中的子序列匹配 shorter 的最大长度.

    >>> matchSubsequence("aabc", "abc")
    [0, 1, 1, 2, 3]
    >>> matchSubsequence("abc", "abcd")
    [0, 1, 2, 3]
    """
    res = [0] * (len(longer) + 1)
    i, j = 0, 0
    while i < len(longer) and j < len(shorter):
        j += longer[i] == shorter[j]
        i += 1
        res[i] = j
    res[i + 1 :] = [j] * (len(longer) - i)
    return res


class Solution:
    def validSequence(self, longer: str, shorter: str) -> List[int]:
        def restore(k: int) -> List[int]:
            """在下标k处可以任意替换字符.返回合法的下标序列."""
            m1 = dp1[k]
            res = []
            i, j = 0, 0
            while len(res) < m1:
                if longer[i] == shorter[j]:
                    res.append(i)
                    j += 1
                i += 1

            f = longer[i] == shorter[j]
            res.append(i)
            i, j = i + 1, j + 1

            while len(res) < len(shorter):
                if longer[i] == shorter[j] or f:
                    f &= longer[i] == shorter[j]
                    res.append(i)
                    j += 1
                i += 1

            return res

        dp1 = matchSubsequence(longer, shorter)
        dp2 = matchSubsequence(longer[::-1], shorter[::-1])[::-1]
        for i in range(len(longer)):
            if dp1[i] + dp2[i + 1] + 1 >= len(shorter):  # [0,i) x [i+1,n)
                return restore(i)
        return []
