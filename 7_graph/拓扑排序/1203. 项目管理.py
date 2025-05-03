# 利用扩展的拓扑排序一次性求出结果
# https://leetcode.cn/problems/sort-items-by-groups-respecting-dependencies/solutions/467233/li-yong-kuo-zhan-de-tuo-bu-pai-xu-yi-ci-xing-qiu-c/

from typing import List


class Solution:
    def sortItems(
        self, n: int, m: int, group: List[int], beforeItems: List[List[int]]
    ) -> List[int]:
        dep = [0] * n
        group_dep = [0] * m
        other_queue = []
        group_queue = [[] for _ in range(m)]
        queue_queue = []
        dependedBy = [[] for _ in range(n)]
        for i, l in enumerate(beforeItems):
            dep[i] = len(l)
            gi = group[i]
            for j in l:
                dependedBy[j].append(i)
                if gi != -1 and group[j] != gi:
                    group_dep[gi] += 1
            if not l:
                (group_queue[gi] if gi != -1 else other_queue).append(i)
        for i, gd in enumerate(group_dep):
            if not gd:
                queue_queue.append(group_queue[i])
        result = []
        while True:
            if queue_queue:
                current_queue = queue_queue.pop()
            elif other_queue:
                current_queue = other_queue
            else:
                if len(result) < n:
                    return []
                else:
                    return result
            while current_queue:
                p = current_queue.pop()
                result.append(p)
                gp = group[p]
                for i in dependedBy[p]:
                    gi = group[i]
                    dep[i] -= 1
                    if not dep[i]:
                        (group_queue[gi] if gi != -1 else other_queue).append(i)
                    if gi != -1 and gi != gp:
                        group_dep[gi] -= 1
                        if not group_dep[gi]:
                            queue_queue.append(group_queue[gi])
