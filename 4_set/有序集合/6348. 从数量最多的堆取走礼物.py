# 给你一个整数数组 gifts ，表示各堆礼物的数量。每一秒，你需要执行以下操作：

# 选择礼物数量最多的那一堆。
# 如果不止一堆都符合礼物数量最多，从中选择任一堆即可。
# 选中的那一堆留下平方根数量的礼物（向下取整），取走其他的礼物。
# 返回在 k 秒后剩下的礼物数量。

# !k很大怎么办?
# 注意到一个数不断开平方，会在 loglogn 次内到达 1
# 如果最大值为0,break即可


from typing import List
from sortedcontainers import SortedList


class Solution:
    def pickGifts(self, gifts: List[int], k: int) -> int:
        sl = SortedList(gifts)
        for _ in range(k):
            max_ = sl.pop()
            sl.add(int(max_**0.5))
            if sl[-1] == 1:
                break
        return sum(sl)


print(Solution().pickGifts([1, 2, 3, 4, 5], 300000000000000000000000))
