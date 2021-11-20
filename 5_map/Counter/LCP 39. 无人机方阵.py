# from collections import Counter
# from typing import List
# from functools import fl

# # 请返回从 source 到 target 最少需要多少架无人机切换灯光颜色。


from collections import Counter
from typing import List


class Solution:
    def minimumSwitchingTimes(self, source: List[List[int]], target: List[List[int]]) -> int:
        c1, c2, res = Counter(sum(source, [])), Counter(sum(target, [])), 0
        for key in set(c1) | set(c2):
            res += abs(c1[key] - c2[key])
        return res // 2


# 注意可选的第2个参数start
# 是一个相加的起始值，相当于从start这个值开始相加。
print(sum([[1, 2], [3, 4]], []))

