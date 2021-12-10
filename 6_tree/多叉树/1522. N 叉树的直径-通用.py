from collections import defaultdict

"""
# Definition for a Node.
class Node:
    def __init__(self, val=None, children=None):
        self.val = val
        self.children = children if children is not None else []
"""


class Solution:
    def diameter(self, root: 'Node') -> int:
        """
        :type root: 'Node'
        :rtype: int
        """
        self.adjMap = defaultdict(list)  # 邻接表
        self.dfs(root)  # 初始化邻接表，建图

        que = [root]  # 随便选一个点作为起点
        visited = set()
        visited.add(root)  # 记忆化
        cur = root  # 全局变量，记下第一次BFS的最后一个点
        while que:
            cur_len = len(que)
            for _ in range(cur_len):
                cur = que.pop(0)
                for nxt in self.adjMap[cur]:
                    if nxt not in visited:
                        visited.add(nxt)
                        que.append(nxt)
        visited.clear()

        que = [cur]  # 第一次BFS最后一个点最为第二次BFS的起点
        visited.add(cur)  # 记忆化
        level = -1  # 波纹法 记录距离
        while que:
            cur_len = len(que)
            level += 1
            for _ in range(cur_len):  # 波纹法
                cur = que.pop(0)
                for nxt in self.adjMap[cur]:
                    if nxt not in visited:
                        visited.add(nxt)
                        que.append(nxt)
        return level

    def dfs(self, rt: 'Node') -> None:  # 初始化邻接表，建图
        if not rt:
            return
        for ch in rt.children:
            self.adjMap[rt].append(ch)
            self.adjMap[ch].append(rt)

            self.dfs(ch)

