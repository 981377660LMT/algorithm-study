# 设计一个数据结构，有效地找到给定子数组的 出现 threshold 次数或次数以上的元素 。


from typing import List
import numpy as np


class MajorityChecker:
    def __init__(self, arr: List[int]):

        self.arr = np.array(arr)

    def query(self, left: int, right: int, threshold: int) -> int:
        # numpy找到最大出现次数的非负数元素

        counter = np.bincount(self.arr[left : right + 1])
        maxPos = np.argmax(counter)
        return int(maxPos) if counter[maxPos] >= threshold else -1
