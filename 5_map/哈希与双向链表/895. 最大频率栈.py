# 895. 最大频率栈
# https://leetcode.cn/problems/maximum-frequency-stack/description/?envType=problem-list-v2&envId=design

from collections import defaultdict


class FreqStack:
    def __init__(self):
        self.stacks = []
        self.counter = defaultdict(int)

    def push(self, val: int) -> None:
        c = self.counter[val]
        if c == len(self.stacks):
            self.stacks.append([val])
        else:
            self.stacks[c].append(val)
        self.counter[val] += 1

    def pop(self) -> int:
        """
        删除并返回堆栈中出现频率最高的元素。
        如果出现频率最高的元素不只一个，则移除并返回最接近栈顶的元素。
        """
        res = self.stacks[-1].pop()
        if not self.stacks[-1]:
            self.stacks.pop()
        self.counter[res] -= 1
        return res


# Your FreqStack object will be instantiated and called as such:
# obj = FreqStack()
# obj.push(val)
# param_2 = obj.pop()
