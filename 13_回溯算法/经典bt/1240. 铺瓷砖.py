# 1 <= n <= 13
# 1 <= m <= 13
# 房子的客厅大小为 n x m，为保持极简的风格，需要使用尽可能少的 正方形 瓷砖来铺盖地面。
# 请你帮设计师计算一下，最少需要用到多少块方形瓷砖？


# 1. 判断bt结束条件(全看过)
# 2. 找到第一个没看过的位置
# 3. 边长从大到小遍历，判断是否满足下一步bt条件
# 4. 模拟填充，回溯
class Solution:
    def tilingRectangle(self, n: int, m: int) -> int:
        def bt(steps: int) -> None:
            if steps >= self.res:
                return

            # 是否每个位置都被覆盖了
            is_all_visited = True
            # 未被覆盖的第一个位置
            cur_row = -1
            cur_col = -1

            for r in range(row):
                for c in range(col):
                    if not visited[r][c]:
                        is_all_visited = False
                        cur_row = r
                        cur_col = c
                        break
                if not is_all_visited:
                    break

            if is_all_visited:  # 该种方案，都被覆盖了，就更新 self.res
                self.res = min(self.res, steps)
                return

            for side in range(min(row - cur_row, col - cur_col), 0, -1):  # 下一步用的砖的边长，从大到小，加速剪枝
                can_put = True

                for r in range(cur_row, cur_row + side):
                    for c in range(cur_col, cur_col + side):
                        if visited[r][c]:
                            can_put = False
                            break
                    if not can_put:
                        break

                if can_put:
                    # 标记
                    for r in range(cur_row, cur_row + side):
                        for c in range(cur_col, cur_col + side):
                            visited[r][c] = True
                    bt(steps + 1)
                    # 回溯
                    for r in range(cur_row, cur_row + side):
                        for c in range(cur_col, cur_col + side):
                            visited[r][c] = False

        row, col = n, m
        if row == col:
            return 1

        visited = [[False for _ in range(col)] for _ in range(row)]
        self.res = row * col  # 最多的情况就是，全用1X1的方砖

        bt(0)
        return self.res


print(Solution().tilingRectangle(5, 8))

