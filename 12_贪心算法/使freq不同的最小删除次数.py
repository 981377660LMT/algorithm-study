# 使字符出现频率不同的最小删除次数
from collections import Counter


class Solution:
    def solve(self, s):
        counter = Counter(s)
        visited = set()
        res = 0
        for f in counter.values():
            if f in visited:
                while f in visited and f > 0:
                    res += 1
                    f -= 1
            visited.add(f)
        return res


print(Solution().solve(s="aabb"))
