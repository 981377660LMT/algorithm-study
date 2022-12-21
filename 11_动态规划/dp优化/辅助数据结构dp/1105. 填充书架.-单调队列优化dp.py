from typing import List, Tuple


# 1105.填充书架.-单调队列优化dp
# 摆放书的顺序与你整理好的顺序相同。
# 以这种方式布置书架，返回书架整体可能的最小高度。
# 1 <= books.length <= 1000
# 1 <= thicknessi <= shelfWidth <= 1000
# 1 <= heighti <= 1000

# https://zhuanlan.zhihu.com/p/374477276
# https://blog.csdn.net/emmmmmmmmm/article/details/83825825
# https://blog.csdn.net/anheku1562/article/details/101308784?spm=1001.2101.3001.6650.1&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-101308784-blog-83825825.pc_relevant_3mothn_strategy_and_data_recovery&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-101308784-blog-83825825.pc_relevant_3mothn_strategy_and_data_recovery&utm_relevant_index=2
# !f[i] = min(f[i], f[j - 1] + max(height[j : i + 1]) 且 sum(width[j : i +1]) < shelfWidth
# !0.只需要每一段的头
# !1.求出新加入的height的影响范围,然后从队列中pop出所有小于这个height的元素
# max可以由滑窗最大值维护


INF = int(1e18)


class Solution:
    def minHeightShelves(self, books: List[List[int]], shelfWidth: int) -> int:
        """O(n) dp[i]表示前i本书的最小高度"""
        n = len(books)
        heights = [0] + [h for _, h in books]
        preSum = [0] * (n + 1)
        dp = [0] * (n + 1)
        queue = MinDeque()
        for i in range(1, n + 1):
            preSum[i] = preSum[i - 1] + books[i - 1][0]
            left = i
            while queue and heights[queue[-1][2]] <= heights[i]:
                left = queue.pop()[1]
            queue.append((dp[left - 1] + heights[i], left, i))
            while queue and preSum[i] - preSum[queue[0][1] - 1] > shelfWidth:
                _, left, right = queue.popleft()
                if left + 1 <= right:
                    queue.appendleft((dp[left] + heights[right], left + 1, right))
            dp[i] = queue.min
        return dp[n]


E = Tuple[int, ...]  # (value, index,...)


class MinDeque:
    __slots__ = ("_left", "_right")

    def __init__(self):
        self._left = _MinStack()
        self._right = _MinStack()

    def append(self, x: "E") -> None:
        self._right.append(x)

    def appendleft(self, x: "E") -> None:
        self._left.append(x)

    def pop(self) -> "E":
        if self._right:
            return self._right.pop()
        tmp = []
        n = len(self._left)
        for _ in range(n):
            tmp.append(self._left.pop())
        half = n // 2
        for i in range(half - 1, -1, -1):
            self._left.append(tmp[i])
        for i in range(half, n):
            self._right.append(tmp[i])
        return self._right.pop()

    def popleft(self) -> "E":
        if self._left:
            return self._left.pop()
        tmp = []
        n = len(self._right)
        for _ in range(n):
            tmp.append(self._right.pop())
        half = n // 2
        for i in range(half - 1, -1, -1):
            self._right.append(tmp[i])
        for i in range(half, n):
            self._left.append(tmp[i])
        return self._left.pop()

    @property
    def min(self) -> int:
        if not self._left:
            return self._right.min
        if not self._right:
            return self._left.min
        return self._left.min if self._left.min < self._right.min else self._right.min

    def __len__(self) -> int:
        return len(self._left) + len(self._right)

    def __getitem__(self, i: int) -> "E":
        n = len(self)
        if i < 0:
            i += n
        if i < 0 or i >= n:
            raise IndexError("deque index out of range")
        if i < len(self._left):
            return self._left[-i - 1]
        return self._right[i - len(self._left)]

    def __repr__(self) -> str:
        sb = []
        for i in range(len(self)):
            sb.append(str(self[i]))
        return f"MinDeque({', '.join(sb)})"


class _MinStack:
    __slots__ = ("_stack", "_mins")

    def __init__(self):
        self._stack = []
        self._mins = []

    def append(self, x: "E") -> None:
        self._stack.append(x)
        if not self._mins or x[0] <= self._mins[-1]:
            self._mins.append(x[0])

    def pop(self) -> "E":
        res = self._stack.pop()
        if res[0] == self._mins[-1]:
            self._mins.pop()
        return res

    def top(self) -> "E":
        return self._stack[-1]

    @property
    def min(self) -> int:
        return self._mins[-1]

    def __len__(self) -> int:
        return len(self._stack)

    def __getitem__(self, i: int) -> "E":
        return self._stack[i]

    def __repr__(self) -> str:
        return f"{self._stack}"


# 6


# # books = [[1,3],[2,4],[3,2]], shelfWidth = 6
# print(Solution().minHeightShelves(books=[[1, 3], [2, 4], [3, 2]], shelfWidth=6))

# 输出：6
# 解释：
# 3 层书架的高度和为 1 + 3 + 2 = 6 。
# 第 2 本书不必放在第一层书架上。
print(
    Solution().minHeightShelves(
        books=[[1, 1], [2, 3], [2, 3], [1, 1], [1, 1], [1, 1], [1, 2]], shelfWidth=4
    )
)
# assert (
#     Solution().minHeightShelves(books=[[9, 9], [5, 4], [3, 1], [1, 5], [7, 3]], shelfWidth=10) == 17
# )

# assert Solution().minHeightShelves([[7, 3], [8, 7], [2, 7], [2, 5]], 10) == 15

# assert Solution().minHeightShelves(books=[[1, 3], [2, 4], [3, 2]], shelfWidth=6) == 4
# # [[1,1],[2,3],[2,3],[1,1],[1,1],[1,1],[1,2]]
# # 4

# assert (
#     Solution().minHeightShelves(
#         books=[[1, 1], [2, 3], [2, 3], [1, 1], [1, 1], [1, 1], [1, 2]], shelfWidth=4
#     )
#     == 6
# )
print(Solution().minHeightShelves(books=[[1, 1] for _ in range(100000)], shelfWidth=4))
