from functools import lru_cache

# 贪心算法是：每次固定字符串最左边的字母 aa 不变，
# 找出距离字符串右侧最近的 aa，把它交换到字符串最右边。
# 这样字符串的头尾字母就相等了。把字符串的头尾去掉，就变成了子问题。
# 把所有子问题的答案加起来就是最少交换次数。


class Solution:
    def minMovesToMakePalindrome1(self, s: str) -> int:
        """递归"""

        @lru_cache(None)
        def dfs(s: str) -> int:
            n = len(s)
            if n <= 2:
                return 0

            # 左边找第一个等于右端点的字符，找不到就从右边找第一个等于左端点的字符
            for i in range(n - 1):
                if s[i] == s[-1]:
                    return i + dfs(s[:i] + s[i + 1 : -1])

            for i in range(n - 1, 0, -1):
                if s[i] == s[0]:
                    return n - 1 - i + dfs(s[1:i] + s[i + 1 :])

            return -1

        return dfs(s)

    def minMovesToMakePalindrome2(self, s: str) -> int:
        """贪心暴力"""
        chars = list(s)
        res = 0
        while chars:
            if chars == chars[::-1]:
                return res
            # 只出现一次的字符
            if chars.count(chars[-1]) == 1:
                chars = chars[::-1]
                continue
            i = chars.index(chars[-1])
            res += i
            chars = chars[:i] + chars[i + 1 : -1]
        return res


print(Solution().minMovesToMakePalindrome1("scpcyxprxxsjyjrww"))
print(Solution().minMovesToMakePalindrome1(s="aabb"))
print(Solution().minMovesToMakePalindrome1(s="letelt"))
