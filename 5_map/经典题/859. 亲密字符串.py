# 只要我们可以通过交换且仅交換一次 s 中的两个字母得到与 goal 相等的结果，就返回 true ；否则返回 false 。


class Solution:
    def buddyStrings(self, s: str, goal: str) -> bool:
        if len(s) != len(goal):
            return False

        # 字符串相同时，需要交换重复的元素
        if s == goal:
            return len(set(s)) < len(s)

        # 挑出不同的字符对,对数只能为2，并且对称，如 (a,b)与(b,a)
        diff = [(a, b) for a, b in zip(s, goal) if a != b]
        return len(diff) == 2 and diff[0] == diff[1][::-1]


print(Solution().buddyStrings('ab', 'ba'))
print(Solution().buddyStrings('ab', 'ab'))
print(Solution().buddyStrings('aa', 'aa'))

