from typing import List
from collections import deque

# shift[i] = [direction, amount]：
# direction 可以为 0 （表示左移）或 1 （表示右移）。
# amount 表示 s 左右移的位数。
# 左移 1 位表示移除 s 的第一个字符，并将该字符插入到 s 的结尾。
# 类似地，右移 1 位表示移除 s 的最后一个字符，并将该字符插入到 s 的开头。


class Solution:
    def __stringShift(self, s: str, shift: List[List[int]]) -> str:
        d = deque(s)
        for direction, amount in shift:
            if direction == 0:
                d.rotate(-amount)
            else:
                d.rotate(amount)
        return ''.join(d)

    # 批量更新
    def stringShift(self, s: str, shift: List[List[int]]) -> str:
        move = sum(-v if k == 0 else v for k, v in shift) % len(s)
        return s[-move:] + s[:-move]


print(Solution()._Solution__stringShift("abc", [[0, 1], [1, 2]]))
print(Solution().stringShift("abc", [[0, 1], [1, 2]]))
