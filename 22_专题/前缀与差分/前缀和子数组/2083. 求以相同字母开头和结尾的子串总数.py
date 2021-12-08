from collections import defaultdict

# 返回 s 中以相同字符开头和结尾的子字符串总数。
# 1788. 最大化花园的美观度.py
class Solution:
    def numberOfSubstrings(self, s: str) -> int:
        indexMap = defaultdict(list)
        for i, char in enumerate(s):
            indexMap[char].append(i)
        return sum(len(indexes) * (len(indexes) + 1) // 2 for indexes in indexMap.values())


print(Solution().numberOfSubstrings(s="abcba"))
