from string import ascii_lowercase

# firstString字符串中从i位置到j位置的子串(包括j位置的字符)和secondString字符串从a位置到b位置的子串(包括b位置字符)相等
# j-a的数值是所有符合前面三个条件的四元组中可能的最小值(j前a后)

# 返回符合上述 4 个条件的四元组的 个数 。
# 总结：单字符就是最优解。直接统计单字符即可。


class Solution:
    def countQuadruples(self, firstString: str, secondString: str) -> int:
        first, last = {}, {}

        for i, char in enumerate(firstString):
            first.setdefault(char, i)

        for i, char in enumerate(secondString):
            last[char] = i

        resMin = int(1e20)
        res = 0
        for char in ascii_lowercase:
            if char in first and char in last:
                cand = first[char] - last[char]
                if cand < resMin:
                    res = 1
                    resMin = cand
                elif cand == resMin:
                    res += 1

        return res


print(Solution().countQuadruples(firstString="abcd", secondString="bccda"))
# 输出：1
# 解释：(0,0,4,4)是唯一符合条件的四元组且其j-a的数值是最小的.

