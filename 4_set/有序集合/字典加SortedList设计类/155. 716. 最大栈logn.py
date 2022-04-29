from sortedcontainers import SortedList


# python 有序字典+有序集合
# 有序字典按 push 顺序维护 <id,x>，有序集合按大小顺序维护 <x, id>。
# 调用 top 方法的时间复杂度为 O(1) ，调用其他方法的时间复杂度为 O(logn) 。
class MaxStack:
    def __init__(self):
        self.id = 0
        self.id_to_value = dict()
        self.sorted_values = SortedList()

    def push(self, x: int) -> None:
        self.id_to_value[self.id] = x
        self.sorted_values.add((x, self.id))
        self.id += 1

    def pop(self) -> int:
        id, value = self.id_to_value.popitem()
        self.sorted_values.remove((value, id))
        return value

    # 此处可以O(1) 只是dict没提供接口
    def top(self) -> int:
        return next(reversed(self.id_to_value.values()))

    def peekMax(self) -> int:
        return self.sorted_values[-1][0]

    # 弹出最大的
    # 如果有多个最大值 则弹出靠近栈顶的
    def popMax(self) -> int:
        value, id = self.sorted_values.pop()
        self.id_to_value.pop(id)
        return value


# Your MaxStack object will be instantiated and called as such:
# obj = MaxStack()
# obj.push(x)
# param_2 = obj.pop()
# param_3 = obj.top()
# param_4 = obj.peekMax()
# param_5 = obj.popMax()
