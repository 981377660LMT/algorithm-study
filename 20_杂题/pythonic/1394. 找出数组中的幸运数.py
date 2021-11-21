# 如果一个整数的出现频次和它的数值大小相等，我们就称这个整数为「幸运数」。
# 如果数组中存在多个幸运数，只需返回 最大 的那个。
# 如果数组中不含幸运数，则返回 -1 。
from typing import List
from collections import Counter


class Solution:
    def findLucky(self, arr: List[int]) -> int:
        return next((num for num, cnt in Counter(arr).most_common() if num == cnt), -1)

