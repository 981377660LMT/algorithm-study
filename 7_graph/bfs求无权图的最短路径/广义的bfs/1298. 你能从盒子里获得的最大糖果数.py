from typing import List

# 时间复杂度：O(N)

# 递进的关系:有盒子=>能开盒子=>开了盒子
# 🔑控制canOpen
class Solution:
    def maxCandies(
        self,
        status: List[int],
        candies: List[int],
        keys: List[List[int]],
        containedBoxes: List[List[int]],
        initialBoxes: List[int],
    ) -> int:
        n = len(status)
        hasBox = [False] * n  # 是否拥有了这个盒子
        canOpen = [status[i] == 1 for i in range(n)]  # 盒子index是否能打开 True/False
        isOpened = [False] * n  # 是否打开了这个盒子

        res = 0
        queue = []  # 都是打开的盒子
        for cur in initialBoxes:  # 初始化
            hasBox[cur] = True  # 拥有了这个盒子
            if canOpen[cur]:  # 如果能打开
                queue.append(cur)  # 进队
                isOpened[cur] = True  # 标记

        while queue:
            cur = queue.pop(0)
            # ---- 1.获取当前盒子中的糖果数
            res += candies[cur]
            # ---- 2.盒子里的钥匙，打开新的盒子
            for key in keys[cur]:
                canOpen[key] = True
                if hasBox[key] and not isOpened[key]:  # 有这个盒子，且未开
                    queue.append(key)
                    isOpened[key] = True
            # ---- 3.继续探索盒子里的子盒子
            for innerBox in containedBoxes[cur]:
                hasBox[innerBox] = True
                if canOpen[innerBox] and not isOpened[innerBox]:
                    queue.append(innerBox)
                    isOpened[innerBox] = True
        return res


print(
    Solution().maxCandies(
        status=[1, 0, 1, 0],
        candies=[7, 5, 4, 100],
        keys=[[], [], [1], []],
        containedBoxes=[[1, 2], [3], [], []],
        initialBoxes=[0],
    )
)

# 输出：16
# 解释：
# 一开始你有盒子 0 。你将获得它里面的 7 个糖果和盒子 1 和 2。
# 盒子 1 目前状态是关闭的，而且你还没有对应它的钥匙。所以你将会打开盒子 2 ，并得到里面的 4 个糖果和盒子 1 的钥匙。
# 在盒子 1 中，你会获得 5 个糖果和盒子 3 ，但是你没法获得盒子 3 的钥匙所以盒子 3 会保持关闭状态。
# 你总共可以获得的糖果数目 = 7 + 4 + 5 = 16 个。

