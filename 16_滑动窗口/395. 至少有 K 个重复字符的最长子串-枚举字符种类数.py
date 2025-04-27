# 395. 至少有 K 个重复字符的最长子串
# https://leetcode.cn/problems/longest-substring-with-at-least-k-repeating-characters/description/
# 给你一个字符串 s 和一个整数 k ，请你找出 s 中的最长子串， 要求该子串中的每一字符出现次数都不少于 k 。
# 返回这一子串的长度。
# 如果不存在这样的子字符串，则返回 0。


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def longestSubstring(self, s: str, k: int) -> int:
        n = len(s)
        if n < k:
            return 0

        res = 0

        for targetUniqueCount in range(1, len(set(s)) + 1):
            counter = [0] * 26
            left, right = 0, 0
            uniqueCount = 0  # 窗口中当前不同字符数
            okCount = 0  # 窗口中符合 ≥k 次的字符数

            while right < n:
                if uniqueCount <= targetUniqueCount:
                    v = ord(s[right]) - ord("a")
                    uniqueCount += counter[v] == 0
                    counter[v] += 1
                    okCount += counter[v] == k
                    right += 1
                else:
                    v = ord(s[left]) - ord("a")
                    okCount -= counter[v] == k
                    counter[v] -= 1
                    uniqueCount -= counter[v] == 0
                    left += 1

                if uniqueCount == targetUniqueCount == okCount:
                    res = max2(res, right - left)

        return res
