class Solution:
    def solve(self, s):
        match = {}
        stack = []
        for i, c in enumerate(s):
            if c == "(":
                stack.append(i)
            elif c == ")":
                match[stack.pop()] = i

        i = 0
        while match.get(i, None) == len(s) - 1 - i:
            i += 1
        return s[i:-i] if i else s


print(Solution().solve("(()())"))
print(Solution().solve(s="(((abc)))"))
print(Solution().solve(s="(((abc)))(d)"))
