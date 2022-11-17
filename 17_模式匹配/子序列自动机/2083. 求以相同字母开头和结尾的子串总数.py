from collections import defaultdict

# 返回 s 中以相同字符开头和结尾的子字符串总数。
# 1788. 最大化花园的美观度.py
class Solution:
    def numberOfSubstrings(self, s: str) -> int:
        indexMap = defaultdict(list)
        for i, char in enumerate(s):
            indexMap[char].append(i)
        return sum(len(indexes) * (len(indexes) + 1) // 2 for indexes in indexMap.values())

    def numberOfSubstrings2(self, s: str) -> int:
        """前缀和"""
        preSum = defaultdict(int)
        res = 0
        for char in s:
            preSum[char] += 1
            res += preSum[char]
        return res


print(Solution().numberOfSubstrings(s="abcba"))
print(Solution().numberOfSubstrings2(s="abcba"))
