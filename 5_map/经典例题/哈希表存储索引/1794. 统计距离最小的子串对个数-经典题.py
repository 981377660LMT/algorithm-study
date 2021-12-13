from string import ascii_lowercase

# firstString字符串中从i位置到j位置的子串(包括j位置的字符)和secondString字符串从a位置到b位置的子串(包括b位置字符)相等
# j-a的数值是所有符合前面三个条件的四元组中可能的最小值(j前a后)

# 返回符合上述 4 个条件的四元组的 个数 。

# 总结：单字符就是最优解。直接统计单字符即可。
class Solution:
    def countQuadruples(self, firstString: str, secondString: str) -> int:
        pos1, pos2 = {}, {}

        # 左边最近
        for i, char in enumerate(firstString):
            if char not in pos1:
                pos1[char] = i

        # 右边最近
        for i, char in enumerate(secondString):
            pos2[char] = i

        minDis = 0x7FFFFFFF
        minDisCount = 0
        for char in ascii_lowercase:
            if char in pos1 and char in pos2:
                dis = pos1[char] - pos2[char]
                if dis < minDis:
                    minDisCount = 1
                    minDis = dis
                elif dis == minDis:
                    minDisCount += 1

        return minDisCount


print(Solution().countQuadruples(firstString="abcd", secondString="bccda"))
# 输出：1
# 解释：(0,0,4,4)是唯一符合条件的四元组且其j-a的数值是最小的.
