# 2067_出现的字符全为k次的子串数-每个字母的前缀和
# 预处理每个位置处每个digit的前缀和


class Solution:
    def equalDigitFrequency(self, s: str) -> int:
        """Given a digit string s, return the number of unique substrings of s where every digit appears the same number of times."""

        def check(left: int, right: int) -> bool:
            """"[left,right] 这一段子串符合题意"""
            diff = set()
            for i in range(10):
                count = preSum[right + 1][i] - preSum[left][i]
                if count > 0:
                    diff.add(count)
                if len(diff) > 1:
                    return False
            return True

        n = len(s)

        # 预处理前缀
        # preSum = [[0] * 10 for _ in range(n + 1)]
        # for i in range(1, n + 1):
        #     preSum[i][ord(s[i - 1]) - ord('0')] += 1
        #     for j in range(10):
        #         preSum[i][j] += preSum[i - 1][j]
        preSum = [[0] * 10]
        for char in s:
            cur = preSum[-1][:]
            cur[ord(char) - ord('0')] += 1
            preSum.append(cur)

        res = set()
        # 枚举所有子串
        for i in range(n):
            for j in range(i, n):
                if check(i, j):
                    res.add(s[i : j + 1])

        return len(res)


print(Solution().equalDigitFrequency(s="1212"))
