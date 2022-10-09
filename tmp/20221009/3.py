from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 和一个机器人，机器人当前有一个空字符串 t 。执行以下操作之一，直到 s 和 t 都变成空字符串：

# 删除字符串 s 的 第一个 字符，并将该字符给机器人。机器人把这个字符添加到 t 的尾部。
# 删除字符串 t 的 最后一个 字符，并将该字符给机器人。机器人将该字符写到纸上。


from typing import List, Tuple


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
        book = (["|"] + list(accumulate(s[:0:-1], min)))[::-1]
        stack, res = [], []
        for ch, lo in zip(s, book):
            stack.append(ch)
            while stack and stack[-1] <= lo:
                res.append(stack.pop())
        return "".join(res)


print(Solution().robotWithString(s="zza"))
print(Solution().robotWithString(s="bac"))
print(Solution().robotWithString(s="bdda"))
print(Solution().robotWithString2(s="bdda"))
