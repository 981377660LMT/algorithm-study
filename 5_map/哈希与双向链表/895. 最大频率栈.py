from collections import defaultdict


class FreqStack:
    def __init__(self):
        self.freq = defaultdict(int)  # 记录每个值出现的频率
        self.group = defaultdict(list)  # 记录每个频率对应的元素
        self.maxfreq = 0

    def push(self, val: int) -> None:
        self.freq[val] += 1
        self.group[self.freq[val]].append(val)
        self.maxfreq = max(self.maxfreq, self.freq[val])

    def pop(self) -> int:
        """
        删除并返回堆栈中出现频率最高的元素。
        如果出现频率最高的元素不只一个，则移除并返回最接近栈顶的元素。
        """
        popped = self.group[self.maxfreq].pop()
        self.freq[popped] -= 1
        if not self.group[self.maxfreq]:
            self.maxfreq -= 1
        return popped


# Your FreqStack object will be instantiated and called as such:
# obj = FreqStack()
# obj.push(val)
# param_2 = obj.pop()
