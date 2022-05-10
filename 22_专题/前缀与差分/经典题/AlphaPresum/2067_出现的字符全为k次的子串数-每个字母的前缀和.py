# 1 <= s.length <= 3 * 104


# 思路：`要看一个字母在子串中的出现次数，我们可以用前缀和O(n)处理，O(1)查询`

# 所求字符串的长度只能为count-26 * count之间,
# 我们遍历字符串,同时记录每个位置的前缀字符个数
# 最后检查前缀和之差是否为count或0(不出现)
class Solution:
    def equalCountSubstrings(self, s: str, count: int) -> int:
        def check(left: int, right: int) -> bool:
            """"[left,right] 这一段子串符合题意"""
            if (right - left + 1) % count != 0:
                return False
            for i in range(26):
                diff = preSum[right + 1][i] - preSum[left][i]
                if diff > 0 and diff != count:
                    return False
            return True

        n = len(s)
        res = 0

        # 预处理前缀
        # preSum = [[0] * 26 for _ in range(n + 1)]
        # for i in range(1, n + 1):
        #     preSum[i][ord(s[i - 1]) - ord('a')] += 1
        #     for j in range(26):
        #         preSum[i][j] += preSum[i - 1][j]
        preSum = [[0] * 26]
        for char in s:
            cur = preSum[-1][:]
            cur[ord(char) - ord('a')] += 1
            preSum.append(cur)

        # 注意所求字符串的长度只能在count-26 * count之间
        for i in range(n):
            for j in range(i + count - 1, min(n, i + 26 * count), count):
                if check(i, j):
                    res += 1
        return res


print(Solution().equalCountSubstrings(s="aaabcbbcc", count=3))

#  "aaa".
#  "bcbbcc".
#  "aaabcbbcc".
