# 给一个字符串 s 和一个字符串 t ，计算 s 中包含 t 中所有字符的子字符串个数。
# 1<=s.length,t.length<=1e5
from collections import Counter

# 1358. 包含所有三种字符的子字符串数目
class Solution:
    def solve(self, s: str, t: str) -> int:
        res, left, need, cur = 0, 0, Counter(t), Counter()
        for char in s:
            cur[char] += 1
            while True:
                next = cur.copy()
                next[s[left]] -= 1
                if next & need == need:
                    cur, left = next, left + 1
                else:
                    break
            if cur & need == need:
                res += left + 1
        return res

    def solve2(self, s: str, t: str) -> int:
        res, left, need, cur = 0, 0, Counter(t), Counter()
        for char in s:
            cur[char] += 1
            while cur & need == need:
                cur[s[left]] -= 1
                left += 1
            res += left
        return res


print(Solution().solve(s="aa", t="a"))
print(Solution().solve(s="acba", t="a"))
print(Solution().solve(s="acba", t="ca"))
print(Solution().solve2(s="aa", t="a"))
print(Solution().solve2(s="acba", t="a"))
print(Solution().solve2(s="acba", t="ca"))
