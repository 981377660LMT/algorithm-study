# 1240. 铺瓷砖
# 1 <= n <= 13
# 1 <= m <= 13
# 房子的客厅大小为 n x m，为保持极简的风格，需要使用尽可能少的 正方形 瓷砖来铺盖地面。
# 请你帮设计师计算一下，最少需要用到多少块正方形瓷砖？
#
#
# !1. bt(count) 回溯函数结束条件:全看过
# !2. 遍历整个矩阵寻找转移点,找到`第一个没看过的位置`
# !3. 边长从大到小遍历，判断是否满足下一步bt条件
# !4. 模拟填充，回溯
#
# 这个问题的研究状况总结
# https://leetcode.cn/problems/tiling-a-rectangle-with-the-fewest-squares/solutions/2316994/zhe-ge-wen-ti-de-yan-jiu-zhuang-kuang-zo-b7jv/


class Solution:
    def tilingRectangle(self, n: int, m: int) -> int:
        if n == m:
            return 1

        self.rect = [[0] * m for _ in range(n)]
        self.res = n * m
        self.backtrack(0)
        return self.res

    def backtrack(self, count: int):
        if count >= self.res:
            return

        # 寻找第一个未覆盖的位置（左上角策略）
        r, c = -1, -1
        for i in range(len(self.rect)):
            for j in range(len(self.rect[0])):
                if self.rect[i][j] == 0:
                    r, c = i, j
                    break
            if r != -1:
                break

        if r == -1:
            self.res = count
            return

        max_size = min(len(self.rect) - r, len(self.rect[0]) - c)
        for size in range(max_size, 0, -1):
            if self.can_place(r, c, size):
                self.place(r, c, size, 1)
                self.backtrack(count + 1)
                self.place(r, c, size, 0)

    def can_place(self, r, c, size):
        # 检查是否可以在(r,c)位置放置size×size的正方形
        if r + size > len(self.rect) or c + size > len(self.rect[0]):
            return False

        for i in range(r, r + size):
            for j in range(c, c + size):
                if self.rect[i][j] == 1:
                    return False
        return True

    def place(self, r, c, size, val):
        # 在(r,c)位置放置或移除size×size的正方形
        for i in range(r, r + size):
            for j in range(c, c + size):
                self.rect[i][j] = val


print(Solution().tilingRectangle(5, 8))
