from graphlib import CycleError, TopologicalSorter


class Solution:
    def canFinish(self, numCourses: int, prerequisites: list[list[int]]) -> bool:
        ts = TopologicalSorter()
        for cur, pre in prerequisites:
            ts.add(cur, pre)
        try:
            ts.prepare()
            return True
        except CycleError:
            return False


print(Solution().canFinish(2, [[1, 0], [0, 1]]))
print(Solution().canFinish(2, [[1, 0]]))
