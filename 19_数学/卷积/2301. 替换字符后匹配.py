from collections import defaultdict
from typing import Any, List


# 给你两个字符串 s 和 sub 。同时给你一个二维字符数组 mappings ，
# 其中 mappings[i] = [oldi, newi] 表示你可以将 sub 中任意数目的 oldi 字符替换为 newi 。sub 中每个字符 不能 被替换超过一次。
# 如果使用 mappings 替换 0 个或者若干个字符，可以将 sub 变成 s 的一个子字符串，请你返回 true，否则返回 false 。
# 一个 子字符串 是字符串中连续非空的字符序列。

# https://leetcode.cn/problems/match-substring-after-replacement/solution/shu-ju-fan-wei-geng-da-de-hua-zen-yao-zu-d9es/

# 令 n=|s|，m=|sub|。对于每种字符 c，
# 定义布尔数组 A[0,…,n−1]，
# 其中 A[i]=1 当且仅当 s[i]=c。
# 再定义布尔数组 B[0,…,m−1]，
# 其中 B[j]=1 当且仅当 sub[j] 不能变成 c。
# 那么对于任意满足 A[i]=1 的下标 i，
# 以及对于任意满足 B[j]=1 的下标 j，
# s[i] 不能跟转化后的 sub[j] 匹配上，
# 所以 sub 串的起始位置不能为 i−j。
# 然后用卷积找出所有不合法的起始位置，
# 再对所有字符取 or。总复杂度 O(σ⋅nlogn)，其中 σ 为字符集大小。


class Solution:
    def matchReplacement(self, s: str, sub: str, mappings: List[List[str]]) -> bool:
        """时间复杂度 O(字符集大小 * n * logn)"""
        n, m = len(s), len(sub)
        ok = set((x, y) for x, y in mappings)
        mp1 = defaultdict(lambda: [0] * n)  # 用于记录每种字符出现的索引
        mp2 = defaultdict(lambda: [1] * m)  # 用于记录哪些索引处的字符不能变成key
        for i, c in enumerate(s):
            mp1[c][i] = 1
        for j, c1 in enumerate(sub):
            for c2 in mp1:
                if c1 == c2 or (c1, c2) in ok:
                    mp2[c2][j] = 0

        # 卷积求出所有不合法的起始位置
        # !reverse一下，i+j 就变成i+(m-j-1)了
        bad = [False] * (n - m + 1)
        for key in mp1:
            nums1, nums2 = mp1[key], mp2[key][::-1]
            conv = convolve(nums1, nums2)
            for i in range(m - 1, n):  # 注意这里的范围
                bad[i - m + 1] |= conv[i] != 0

        return not all(bad)


import numpy as np


def convolve(nums1: Any, nums2: Any) -> "np.ndarray":
    """fft求卷积"""
    n, m = len(nums1), len(nums2)
    ph = 1 << (n + m - 2).bit_length()
    T = np.fft.rfft(nums1, ph) * np.fft.rfft(nums2, ph)
    res = np.fft.irfft(T, ph)[: n + m - 1]
    return np.rint(res).astype(np.int64)


print(
    Solution().matchReplacement(
        s="fool3e7bar" * 1000, sub="leet" * 1000, mappings=[["e", "3"], ["t", "7"], ["t", "8"]]
    )
)
# s = "fooleetbar", sub = "f00l", mappings = [["o","0"]]
print(Solution().matchReplacement(s="fooleetbar", sub="f00l", mappings=[["o", "0"]]))
