# 输入：boxes = "110"
# 输出：[1,1,3]
# 解释：每个盒子对应的最小操作数如下：
# 1) 第 1 个盒子：将一个小球从第 2 个盒子移动到第 1 个盒子，需要 1 步操作。
# 2) 第 2 个盒子：将一个小球从第 1 个盒子移动到第 2 个盒子，需要 1 步操作。
# 3) 第 3 个盒子：将一个小球从第 1 个盒子移动到第 3 个盒子，需要 2 步操作。将一个小球从第 2 个盒子移动到第 3 个盒子，需要 1 步操作。共计 3 步操作。

# 其中 answer[i] 是将所有小球移动到第 i 个盒子所需的 最小 操作数。
# 每次只能相邻移动1个球
# 1 <= n <= 2000

from typing import List


class Solution:
    # 暴力O(n^2)
    def minOperations2(self, boxes: str) -> List[int]:
        res = []
        for i in range(len(boxes)):
            count = 0
            for j in range(len(boxes)):
                if boxes[j] == '1':
                    count += abs(j - i)
            res.append(count)

        return res

    def minOperations(self, boxes: str) -> List[int]:
        res = [0] * len(boxes)
        notEmpty, runningSum = 0, 0

        for i, box in enumerate(boxes):
            res[i] += runningSum
            if box == '1':
                notEmpty += 1
            runningSum += notEmpty

        notEmpty, runningSum = 0, 0

        for i, box in reversed(list(enumerate(boxes))):
            res[i] += runningSum
            if box == '1':
                notEmpty += 1
            runningSum += notEmpty

        return res


print(Solution().minOperations('110'))
