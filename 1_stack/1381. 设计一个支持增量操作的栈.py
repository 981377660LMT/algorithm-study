# 1381. 设计一个支持增量操作的栈
# https://leetcode.cn/problems/design-a-stack-with-increment-operation/description/
# 延迟操作: 记录变化到差分哈希表，po p出时再执行增加操作, 并向前累加delta


class CustomStack:

    def __init__(self, maxSize: int):
        self.stack = []
        self.incrementals = []
        self.maxSize = maxSize

    def push(self, x: int) -> None:
        if len(self.stack) == self.maxSize:
            return
        self.stack.append(x)
        self.incrementals.append(0)

    def pop(self) -> int:
        if not self.stack:
            return -1
        added = self.incrementals.pop()
        if self.incrementals:
            self.incrementals[-1] += added
        return self.stack.pop() + added

    def increment(self, k: int, val: int) -> None:
        """栈底的 k 个元素的值都增加 val 。如果栈中元素总数小于 k ，则栈中的所有元素都增加 val."""
        if not self.stack:
            return
        bottomK = min(k, len(self.stack))
        self.incrementals[bottomK - 1] += val
