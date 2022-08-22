"""
最大回文数字
回文的组成:偶数部分+一个奇数?(可选)+偶数部分
如果有奇数个 那么拆成1个奇数+偶数部分
如果有偶数个 那么拆成偶数部分

返回答案时注意 or '0' 兜底
"""


from collections import Counter


class Solution:
    def largestPalindromic(self, num: str) -> str:
        """请你找出能够使用 num 中数字形成的 最大回文 整数"""
        counter = Counter(num)
        odds = []
        evens = []

        for i in range(9, -1, -1):
            s = str(i)
            count = counter[s]
            if count == 0:
                continue
            if count & 1:
                odds.append(s)
                half = (count - 1) // 2
                if half:
                    evens.append(s * half)
            else:
                half = count // 2
                evens.append(s * half)

        # 前导0的情况
        if len(evens) == 1 and evens[0][0] == "0":
            evens = []

        res1 = "".join(evens)
        cand = res1 + (str(odds[0]) if odds else "") + res1[::-1]
        return cand or "0"  # !注意有些回文题需要这样保底


# print(Solution().largestPalindromic(num="444947137"))
# print(Solution().largestPalindromic(num="010"))
# print(Solution().largestPalindromic(num="00009"))
print(Solution().largestPalindromic(num="00001105"))
# "1005001"
# print(Solution().largestPalindromic(num="00011"))
# "10001"
# print(Solution().largestPalindromic(num="12345"))
# print(Solution().largestPalindromic(num="9988099"))
# # hidden
# print(Solution().largestPalindromic(num="0009"))
# print(Solution().largestPalindromic(num="0008765"))
# print(Solution().largestPalindromic(num="00"))
