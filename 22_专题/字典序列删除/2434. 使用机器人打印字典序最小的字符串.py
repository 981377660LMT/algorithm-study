# 字典序最小的出栈序列

from itertools import accumulate
from collections import Counter


# 给你一个字符串 s 和一个机器人，机器人当前有一个空字符串 t 。
# 执行以下操作之一，直到 s 和 t 都变成空字符串：
# 删除字符串 s 的 第一个 字符，并将该字符给机器人。机器人把这个字符添加到 t 的尾部。
# 删除字符串 t 的 最后一个 字符，并将该字符给机器人。机器人将该字符写到纸上。


class Solution:
    def robotWithString(self, s: str) -> str:
        stack, remain, res = [], Counter(s), []
        min_ = ord("a")
        for char in s:
            stack.append(char)
            remain[char] -= 1
            while min_ < ord("z") and remain[chr(min_)] == 0:
                min_ += 1
            while stack and ord(stack[-1]) <= min_:
                res.append(stack.pop())
        return "".join(res)

    def robotWithString2(self, s: str) -> str:
        """后缀最小值的写法"""
        BIG = chr(200)
        sufMin = ([BIG] + list(accumulate(s[::-1], min)))[::-1]
        stack, res = [], []
        for i, char in enumerate(s):
            stack.append(char)
            while stack and stack[-1] <= sufMin[i + 1]:
                res.append(stack.pop())
        return "".join(res)


print(Solution().robotWithString(s="zza"))
print(Solution().robotWithString(s="bac"))
print(Solution().robotWithString(s="bdda"))
print(Solution().robotWithString2(s="bdda"))
