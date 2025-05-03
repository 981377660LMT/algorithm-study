from typing import List
from collections import defaultdict


# 有 n 个项目，每个项目或者不属于任何小组，或者属于 m 个小组之一
# group[i] 表示第 i 个项目所属的小组
# 如果第 i 个项目不属于任何小组，则 group[i] 等于 -1
# 同一小组的项目，在res中彼此相邻。  =>  同组彼此相邻表示这是 `defaultdict(list)形式` 的结构
# beforeItems[i] 表示在进行第 i 个项目前（位于第 i 个项目左侧）应该完成的所有项目。
class Solution:
    def sortItems(
        self, n: int, m: int, group: List[int], beforeItems: List[List[int]]
    ) -> List[int]:
        # 0.返回拓扑排序结果
        def get_top_order(graph, indegree):
            top_order = []
            stack = [i for i in range(len(graph)) if indegree[i] == 0]
            while stack:
                cur = stack.pop()
                top_order.append(cur)
                for next in graph[cur]:
                    indegree[next] -= 1
                    if indegree[next] == 0:
                        stack.append(next)
            return top_order if len(top_order) == len(graph) else []

        # 1.给没有归属的项目分配所属group的id
        for g in range(len(group)):
            if group[g] == -1:
                group[g] = m
                m += 1

        # 2. 给每个item与group建图
        graph_items = [[] for _ in range(n)]
        indegree_items = [0] * n
        graph_groups = [[] for _ in range(m)]
        indegree_groups = [0] * m
        for u in range(n):
            for v in beforeItems[u]:
                graph_items[v].append(u)
                indegree_items[u] += 1
                if group[u] != group[v]:
                    # v的小组必须在u前面
                    graph_groups[group[v]].append(group[u])
                    indegree_groups[group[u]] += 1

        # 3.获得两个图的拓扑序
        item_order = get_top_order(graph_items, indegree_items)
        group_order = get_top_order(graph_groups, indegree_groups)
        if not item_order or not group_order:
            return []

        # 4.把项目先在组内排好序
        order_within_group = defaultdict(list)
        for v in item_order:
            order_within_group[group[v]].append(v)

        # 5.把排好序的组按顺序组合在一起
        res = []
        for g in group_order:
            res.extend(order_within_group[g])
        return res


print(
    Solution().sortItems(
        n=8,
        m=2,
        group=[-1, -1, 1, 0, 0, 1, 0, -1],
        beforeItems=[[], [6], [5], [6], [3, 6], [], [], []],
    )
)
