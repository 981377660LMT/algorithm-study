from typing import List
from sortedcontainers import SortedDict


class SummaryRanges:
    def __init__(self):
        # self.LR : 区间左端点的右端点
        # self.RL : 区间右端点的左端点
        self.RL = SortedDict()
        self.LR = SortedDict()

    def addNum(self, val: int) -> None:
        # --------先判断val是否已经在某个区间内
        if self.RL:
            pos = self.RL.bisect_left(val)  # 大于等于val的第一个R端点
            if pos < len(self.RL):
                cur_section_L = self.RL.values()[pos]  # 当前区间的左端点
                if cur_section_L <= val:
                    return

        # --------若val未出现过，分情况讨论val-1和val+1是否出现过
        l = (val - 1) in self.RL
        r = (val + 1) in self.LR
        # --（1）若两个邻居都存在, [区间1_L, val - 1], val, [val + 1. 区间2_R].
        if l and r:
            section_1_L = self.RL[val - 1]
            section_2_R = self.LR[val + 1]
            self.LR[section_1_L] = section_2_R
            self.RL[section_2_R] = section_1_L
            del self.RL[val - 1]
            del self.LR[val + 1]
        # --（2）左邻居在，有邻居无， [区间1_L, val - 1], val,
        elif l and not r:
            section_1_L = self.RL[val - 1]
            self.LR[section_1_L] = val
            self.RL[val] = section_1_L
            del self.RL[val - 1]
        # --（3）右邻居在，左邻居不在，val, [val + 1. 区间2_R]
        elif not l and r:
            section_2_R = self.LR[val + 1]
            self.LR[val] = section_2_R
            self.RL[section_2_R] = val
            del self.LR[val + 1]
        # --左右邻居都不相连。val自己成一个区间
        else:
            self.LR[val] = val
            self.RL[val] = val

    def getIntervals(self) -> List[List[int]]:
        res = []
        for l, r in self.LR.items():
            res.append([l, r])
        return res


# Your SummaryRanges object will be instantiated and called as such:
# obj = SummaryRanges()
# obj.addNum(val)
# param_2 = obj.getIntervals()
