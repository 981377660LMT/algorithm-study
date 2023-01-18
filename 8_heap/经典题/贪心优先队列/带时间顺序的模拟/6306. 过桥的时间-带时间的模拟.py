# 模拟 过桥的时间


# 共有 k 位工人计划将 n 个箱子从旧仓库移动到新仓库。给你两个整数 n 和 k，
# 以及一个二维整数数组 time ，数组的大小为 k x 4 ，
# 其中 time[i] = [leftToRighti, pickOldi, rightToLefti, putNewi] 。

# 一条河将两座仓库分隔，只能通过一座桥通行。
# 旧仓库位于河的右岸，新仓库在河的左岸。
# 开始时，所有 k 位工人都在桥的左侧等待。
# 为了移动这些箱子，第 i 位工人（下标从 0 开始）可以：

# 从左岸（新仓库）跨过桥到右岸（旧仓库），用时 leftToRighti 分钟。
# 从旧仓库选择一个箱子，并返回到桥边，用时 pickOldi 分钟。不同工人可以同时搬起所选的箱子。
# 从右岸（旧仓库）跨过桥到左岸（新仓库），用时 rightToLefti 分钟。
# 将箱子放入新仓库，并返回到桥边，用时 putNewi 分钟。不同工人可以同时放下所选的箱子。
# 如果满足下面任一条件，则认为工人 i 的 效率低于 工人 j ：

# leftToRighti + rightToLefti > leftToRightj + rightToLeftj
# leftToRighti + rightToLefti == leftToRightj + rightToLeftj 且 i > j
# 工人通过桥时需要遵循以下规则：

# 如果工人 x 到达桥边时，工人 y 正在过桥，那么工人 x 需要在桥边等待。
# 如果没有正在过桥的工人，那么在桥右边等待的工人可以先过桥。如果同时有多个工人在右边等待，那么 效率最低 的工人会先过桥。
# 如果没有正在过桥的工人，且桥右边也没有在等待的工人，同时旧仓库还剩下至少一个箱子需要搬运，此时在桥左边的工人可以过桥。如果同时有多个工人在左边等待，那么 效率最低 的工人会先过桥。
# 所有 n 个盒子都需要放入新仓库，请你返回最后一个搬运箱子的工人 `到达河左岸` 的时间。

# https://leetcode.cn/problems/time-to-cross-a-bridge/solution/by-kpole-8d8x/
# !模拟一个过程，而这个过程一般是按照 `时间顺序` 去执行:
# !0.分析状态与每种状态的优先级,用SortedList存储
# !1.弄清模拟的结束条件 (while ...)
# !2.每次while循环处理中的事件 : while一次性加入所有,if来处理不同event
#    !注意如果没有event要处理,则需要更新时间到下一次状态变化的时间


from typing import List
from sortedcontainers import SortedList

INF = int(1e18)


class Solution:
    def findCrossingTime(self, n: int, k: int, time: List[List[int]]) -> int:
        rightFinish = SortedList()  # (finishTime, id)
        leftFinish = SortedList()  # (finishTime, id)
        leftWait = SortedList(key=lambda x: (-x[0], -x[1]))  # (leftToRight + rightToLeft, id)
        rightWait = SortedList(key=lambda x: (-x[0], -x[1]))  # (leftToRight + rightToLeft, id)
        for i, (leftToRight, pickOld, rightToLeft, putNew) in enumerate(time):
            leftWait.add((leftToRight + rightToLeft, i))

        remain, curTime = n, 0
        while remain > 0 or rightFinish or rightWait:  # !当旧仓库还有货物或者右边还有人要回来时
            while leftFinish and leftFinish[0][0] <= curTime:
                _, id = leftFinish.pop(0)
                leftToRight, _, rightToLeft, _ = time[id]
                leftWait.add((leftToRight + rightToLeft, id))
            while rightFinish and rightFinish[0][0] <= curTime:
                _, id = rightFinish.pop(0)
                leftToRight, _, rightToLeft, _ = time[id]
                rightWait.add((leftToRight + rightToLeft, id))

            if rightWait:  # 右侧有等待的工人
                _, id = rightWait.pop(0)
                *_, rightToLeft, putNew = time[id]
                curTime += rightToLeft
                leftFinish.add((curTime + putNew, id))
            elif leftWait and remain > 0:  # 旧仓库还有货物，并且左侧有等待的工人，就出发搬一个回来
                _, id = leftWait.pop(0)
                leftToRight, pickOld, *_ = time[id]
                curTime += leftToRight
                rightFinish.add((curTime + pickOld, id))
                remain -= 1
            else:  # !此时没有人需要过桥，时间应该过渡到第一个处于「放下」或者「搬起」状态的工人切换状态的时刻。
                next1 = leftFinish[0][0] if leftFinish else INF
                next2 = rightFinish[0][0] if rightFinish else INF
                curTime = min(next1, next2)
        return curTime


# print(Solution().findCrossingTime(n=1, k=3, time=[[1, 1, 2, 1], [1, 1, 3, 1], [1, 1, 4, 1]]))
print(Solution().findCrossingTime(n=3, k=2, time=[[1, 9, 1, 8], [10, 10, 10, 10]]))
# 10000
# 1
# [[1000,1000,1000,1000]]
print(Solution().findCrossingTime(n=10000, k=1, time=[[1000, 1000, 1000, 1000]]))
