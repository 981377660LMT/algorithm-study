from typing import List

# 请你返回学完全部课程所需的最少学期数。
# 如果没有办法做到学完全部这些课程的话，就返回 -1。
class Solution:
    def minimumSemesters(self, n: int, relations: List[List[int]]) -> int:
        indegree = [0 for _ in range(n + 1)]
        adjlist = [[] for _ in range(n + 1)]
        for u, v in relations:
            indegree[v] += 1
            adjlist[u].append(v)

        res = 0
        visited = 0  # 完成的课程数
        cur_queue = [id for id in range(1, n + 1) if indegree[id] == 0]

        while cur_queue:  # 这个学期有课上
            visited += len(cur_queue)  # 这个学期又ac了这么多课
            next_queue = []  # 下个学期可以上什么
            for cur in cur_queue:
                for next in adjlist[cur]:  # 后续的课程里
                    indegree[next] -= 1  # x上完了，y的前导课程少了1门
                    if indegree[next] == 0:  # 如果y的前导课程都上完了
                        next_queue.append(next)  # 下个学期就可以上y了
            cur_queue = next_queue  # 更新，为下个term做准备
            res += 1  # 学期数+1

        return res if visited == n else -1


print(Solution().minimumSemesters(3, [[1, 3], [2, 3]]))

