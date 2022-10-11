# 1 <= n <= 13
# 1 <= m <= 13
# 房子的客厅大小为 n x m，为保持极简的风格，需要使用尽可能少的 正方形 瓷砖来铺盖地面。
# 请你帮设计师计算一下，最少需要用到多少块正方形瓷砖？


# !1. bt(count) 回溯函数结束条件:全看过
# !2. 遍历整个矩阵寻找转移点,找到`第一个没看过的位置`
# !3. 边长从大到小遍历，判断是否满足下一步bt条件
# !4. 模拟填充，回溯
class Solution:
    def tilingRectangle(self, n: int, m: int) -> int:
        def bt(need: int) -> None:
            if need >= self.res:
                return

            allVisited = True  # 是否每个位置都被覆盖了
            startRow, startCol = -1, -1  # 未被覆盖的第一个位置
            for r in range(ROW):
                for c in range(COL):
                    if not visited[r][c]:
                        allVisited = False
                        startRow, startCol = r, c
                        break
                if not allVisited:
                    break

            if allVisited:
                self.res = min(self.res, need)  # 所有瓷砖都被覆盖了就更新res
                return

            for side in range(min(ROW - startRow, COL - startCol), 0, -1):  # !下一步用的砖的边长，从大到小加速剪枝
                canPut = True
                for r in range(startRow, startRow + side):
                    for c in range(startCol, startCol + side):
                        if visited[r][c]:
                            canPut = False
                            break
                    if not canPut:
                        break

                if canPut:
                    # 标记
                    for r in range(startRow, startRow + side):
                        for c in range(startCol, startCol + side):
                            visited[r][c] = True
                    bt(need + 1)
                    # 回溯
                    for r in range(startRow, startRow + side):
                        for c in range(startCol, startCol + side):
                            visited[r][c] = False

        ROW, COL = n, m
        if ROW == COL:  # 如果是正方形直接返回1
            return 1
        visited = [[False for _ in range(COL)] for _ in range(ROW)]
        self.res = ROW * COL  # 最多的情况就是全用1*1的方砖
        bt(0)
        return self.res


print(Solution().tilingRectangle(5, 8))
