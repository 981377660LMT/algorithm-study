from collections import defaultdict


class Solution:
    def uniqueLetterString(self, s: str) -> int:
        n = len(s)
        indexMap = defaultdict(list)
        for i, char in enumerate(s):
            indexMap[char].append(i)

        res = 0
        # 我们枚举每个字符, 找出不含该字符的区间, 用M减去不含该字符的区间所能产生的子字符串就是单个字符对答案的贡献
        for indexes in indexMap.values():
            res += n * (n + 1) // 2  # 字符串子串数
            # 减去不含该字符的区间
            indexes = [-1] + indexes + [n]
            for pre, cur in zip(indexes, indexes[1:]):
                count = cur - pre
                res -= count * (count - 1) // 2

        return res


print(Solution().uniqueLetterString(s="good"))
print(Solution().uniqueLetterString(s="test"))
# 字符串test
# 从0元素开始，1，2，3，3 总计9种
# 从1元素开始，1，2，3 总计6种
# 从2元素开始，1，2 总计3种
# 从3元素开始，1 总计1种
# 结果为，1 + 3 + 6 + 9 = 19种

# 具体的来说,对于test字符串,子字符串的个数是(4+3+2+1) = 10
# t的贡献是10 - (2+1) = 7 (区间“es”中没有t)
# e的贡献是10 - (1) - (2+1) = 6 (区间"t"和"st"中没有e)
# s的贡献是10 - (2+1) - (1) = 6 (区间"te"和"t"中没有s)
# 总计7+6+6=19
