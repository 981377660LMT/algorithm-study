# 3485. 删除元素后 K 个字符串的最长公共前缀
# https://leetcode.cn/problems/longest-common-prefix-of-k-strings-after-removal/description/
#
# 给你一个字符串数组 words 和一个整数 k。
# 对于范围 [0, words.length - 1] 中的每个下标 i，在移除第 i 个元素后的剩余数组中，找到任意 k 个字符串（k 个下标 互不相同）的 最长公共前缀 的 长度。
# 返回一个数组 answer，其中 answer[i] 是 i 个元素的答案。如果移除第 i 个元素后，数组中的字符串少于 k 个，answer[i] 为 0。
#
# 利用 lcp 的性质.
# 把字符串排序后，有着相同前缀的字符串可以聚在一起。
# !有序子数组的 LCP，等于子数组第一个字符串和最后一个字符串的 LCP。
# 设所有长为 k 的子数组中，LCP 的最大长度为 mx，次大长度为 mx2。设 mx 对应的子数组为 [mxI,mxI+k−1]
# 如果删除的字符串不在 [mxI,mxI+k−1] 中，答案是 mx
# 如果删除的字符串在 [mxI,mxI+k−1] 中，答案是 mx2

from typing import List


def lcp(s1: str, s2: str) -> int:
    res = 0
    for a, b in zip(s1, s2):
        if a != b:
            break
        res += 1
    return res


def longestCommonPrefix(words: List[str], k: int) -> List[int]:
    n = len(words)
    if k >= n:
        return [0] * n
    order = sorted(range(n), key=lambda i: words[i])

    # 计算最大 LCP 长度和次大 LCP 长度，同时记录最大 LCP 来自哪里
    max1, max2 = -1, -1
    maxI = -1
    for i in range(n - k + 1):
        lcp_ = lcp(words[order[i]], words[order[i + k - 1]])
        if lcp_ > max1:
            max1, max2 = lcp_, max1
            maxI = i
        elif lcp_ > max2:
            max2 = lcp_

    res = [max1] * n
    for i in order[maxI : maxI + k]:
        res[i] = max2
    return res
