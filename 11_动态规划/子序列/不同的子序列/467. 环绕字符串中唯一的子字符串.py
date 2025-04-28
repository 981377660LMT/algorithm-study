# 467. 环绕字符串中唯一的子字符串
# https://leetcode.cn/problems/unique-substrings-in-wraparound-string/description/


from collections import defaultdict


class Solution:
    def findSubstringInWraproundString(self, s: str) -> int:
        """
        动态规划：dp[c] 记录所有以字符 c 结尾的、符合环绕规律的最长子串长度。
        遍历 s，用 cur_len 追踪以 s[i] 为结尾的当前有效子串长度：
          - 如果 s[i] 与 s[i-1] 是环绕连续（(ord(s[i]) - ord(s[i-1])) % 26 == 1），cur_len += 1
          - 否则 cur_len = 1
        然后更新 dp[s[i]] = max(dp[s[i]], cur_len)。
        最终不同子串的数量是 sum(dp.values())，因为对于每个结尾字符 c，
        长度为 L 的子串中包含了 L 个不同的以 c 结尾的子串，且这些子串之间互不重叠。
        时间 O(n)，空间 O(1)（26 个字母的常数级数组）。
        """
        dp = defaultdict(int)
        k = 0
        for i, c in enumerate(s):
            if i > 0 and (ord(c) - ord(s[i - 1])) % 26 == 1:  # 字符之差为 1 或 -25
                k += 1
            else:
                k = 1
            dp[c] = max(dp[c], k)  # max去重
        return sum(dp.values())


if __name__ == "__main__":
    sol = Solution()
    print(sol.findSubstringInWraproundString("a"))  # 输出 1: {"a"}
    print(sol.findSubstringInWraproundString("cac"))  # 输出 2: {"a","c"}
    print(sol.findSubstringInWraproundString("zab"))  # 输出 6: {"z","a","b","za","ab","zab"}
