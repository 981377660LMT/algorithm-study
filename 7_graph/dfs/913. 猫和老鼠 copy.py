from typing import List
from functools import lru_cache


class Solution:
    def catMouseGame(self, graph: List[List[int]]) -> int:
        @lru_cache(None)
        def search(t, x, y):
            # 如果总步数达到2*N,意味着猫和老鼠各走了N步，但是老鼠到达洞的步数最多只有N,如果过了N步老鼠还没有到达洞，并且猫也没有抓住老鼠，那么就是平局
            if t == int(2 * len(graph)):
                return 0
            # 老鼠和猫在一个位置，猫赢
            if x == y:
                return 2
            # 老鼠进0，老鼠赢
            if x == 0:
                return 1
            if t % 2 == 0:
                # 下一步是老鼠，如果老鼠能在任何一个下一步中赢 返回老鼠赢。（返回1）
                if any(search(t + 1, x_nxt, y) == 1 for x_nxt in graph[x]):
                    return 1
                # 如果任何一个老鼠的下一步返回的是平局，返回平局（0）
                if any(search(t + 1, x_nxt, y) == 0 for x_nxt in graph[x]):
                    return 0
                # 如果前两者都不是，那么猫赢，意味着返回值全是2，也就是不论老鼠走哪一步，猫都赢
                return 2
            else:
                # 下一步是猫，如果猫能在任何一个下一步中赢 返回猫赢。（返回2）
                if any(search(t + 1, x, y_nxt) == 2 for y_nxt in graph[y] if y_nxt != 0):
                    return 2
                # 如果任何一个猫的下一步返回的是平局，返回平局（0）
                if any(search(t + 1, x, y_nxt) == 0 for y_nxt in graph[y] if y_nxt != 0):
                    return 0
                # 如果前两者都不是，那么老鼠赢，意味着返回值全是1，也就是不论猫走哪一步，老鼠都赢
                return 1

        return search(0, 1, 2)

