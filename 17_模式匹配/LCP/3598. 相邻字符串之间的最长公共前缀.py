# 3598. 相邻字符串之间的最长公共前缀
# https://leetcode.cn/problems/longest-common-prefix-between-adjacent-strings-after-removals/description/
# 给你一个字符串数组 words，对于范围 [0, words.length - 1] 内的每个下标 i，执行以下步骤：
#
# 从 words 数组中移除下标 i 处的元素。
# 计算修改后的数组中所有 相邻对 之间的 最长公共前缀 的长度。
# 返回一个数组 answer，其中 answer[i] 是移除下标 i 后，相邻对之间最长公共前缀的长度。如果 不存在 相邻对，或者 不存在 公共前缀，则 answer[i] 应为 0。
#
# 字符串的前缀是从字符串的开头开始延伸到任意位置的子字符串。


from typing import List


def lcp(s1: str, s2: str) -> int:
    res = 0
    for a, b in zip(s1, s2):
        if a != b:
            break
        res += 1
    return res


class Solution:
    def longestCommonPrefix(self, words: List[str]) -> List[int]:
        n = len(words)

        # !计算最大 LCP 长度和次大 LCP 长度，同时记录最大 LCP 来自哪里
        max1, max2 = -1, -1
        maxI = -1
        for i in range(n - 1):
            lcp_ = lcp(words[i], words[i + 1])
            if lcp_ > max1:
                max1, max2 = lcp_, max1
                maxI = i
            elif lcp_ > max2:
                max2 = lcp_

        res = [0] * n
        for i in range(n):
            lcp_ = lcp(words[i - 1], words[i + 1]) if 0 < i < n - 1 else 0
            # 最大lcp未被破坏
            if i != maxI and i != maxI + 1:
                res[i] = max(max1, lcp_)
            else:
                res[i] = max(max2, lcp_)

        return res
