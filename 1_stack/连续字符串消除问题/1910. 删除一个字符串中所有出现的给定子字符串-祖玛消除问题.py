# 祖玛消除问题
# 1. 正则/2. 栈
# 1 <= s.length <= 1000
class Solution:
    def removeOccurrences(self, s: str, part: str) -> str:
        stack = []
        for char in s:
            stack.append(char)
            if ''.join(stack[-len(part) :]) == part:
                for _ in range(len(part)):
                    stack.pop()

        return ''.join(stack)


print(Solution().removeOccurrences(s="daabcbaabcbc", part="abc"))
# 输出："dab"
# 解释：以下操作按顺序执行：
# - s = "daabcbaabcbc" ，删除下标从 2 开始的 "abc" ，得到 s = "dabaabcbc" 。
# - s = "dabaabcbc" ，删除下标从 4 开始的 "abc" ，得到 s = "dababc" 。
# - s = "dababc" ，删除下标从 3 开始的 "abc" ，得到 s = "dab" 。
# 此时 s 中不再含有子字符串 "abc" 。
