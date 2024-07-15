# 100352. 交换后字典序最小的字符串
# https://leetcode.cn/problems/lexicographically-smallest-string-after-a-swap/description/
# 给你一个仅由数字组成的字符串 s，在最多交换一次 相邻 且具有相同 奇偶性 的数字后，返回可以得到的
# 字典序最小的字符串
# !由于最多只能交换一次，在需要交换的情况下，显然是越早交换越好
#
# 交换任意两位
# 670. 最大交换
# https://leetcode.cn/problems/maximum-swap/description/


class Solution:
    def getSmallestString(self, s: str) -> str:
        sb = list(s)
        for i in range(1, len(sb)):
            a, b = sb[i - 1], sb[i]
            if a > b and ord(a) % 2 == ord(b) % 2:
                sb[i - 1], sb[i] = sb[i], sb[i - 1]
                break
        return "".join(sb)
