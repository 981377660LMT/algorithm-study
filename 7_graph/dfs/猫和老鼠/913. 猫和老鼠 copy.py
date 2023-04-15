# 猫和老鼠


from typing import List
from functools import lru_cache


class Solution:
    def catMouseGame(self, graph: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(turn: int, mouse: int, cat: int) -> int:
            # if turn == 2 * n * n:  # 2*n 不足以判断平局 2*n*n 才可以 但是2*n*n会超时 所以用100估计边界
            #     return 0
            if turn == 100:
                return 0
            # 老鼠和猫在一个位置，猫赢
            if mouse == cat:
                return 2
            # 老鼠进0，老鼠赢
            if mouse == 0:
                return 1

            if not turn & 1:
                # 下一步是老鼠，如果老鼠能在任何一个下一步中赢 返回老鼠赢。（返回1）
                if any(dfs(turn + 1, next, cat) == 1 for next in graph[mouse]):
                    return 1
                # 如果任何一个老鼠的下一步返回的是平局，返回平局（0）
                if any(dfs(turn + 1, next, cat) == 0 for next in graph[mouse]):
                    return 0
                # 如果前两者都不是，那么猫赢，意味着返回值全是2，也就是不论老鼠走哪一步，猫都赢
                return 2
            else:
                # 猫无法移动到洞中
                # 下一步是猫，如果猫能在任何一个下一步中赢 返回猫赢。（返回2）
                if any(dfs(turn + 1, mouse, next) == 2 for next in graph[cat] if next != 0):
                    return 2
                # 如果任何一个猫的下一步返回的是平局，返回平局（0）
                if any(dfs(turn + 1, mouse, next) == 0 for next in graph[cat] if next != 0):
                    return 0
                # 如果前两者都不是，那么老鼠赢，意味着返回值全是1，也就是不论猫走哪一步，老鼠都赢
                return 1

        n = len(graph)
        return dfs(0, 1, 2)
