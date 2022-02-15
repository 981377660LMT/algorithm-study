# https://www.cnblogs.com/wenruo/p/5264235.html

# O(n^3)
from typing import List


INF = int(1e20)


class KM:
    """KM算法求二分图最大权完美匹配"""

    def __init__(self, adjMatrix: List[List[int]]):
        max_ = max(len(adjMatrix), len(adjMatrix[0]))
        self._match = [-1] * max_  # 记录每个女生匹配到的男生 如果没有则为-1
        self._graph = adjMatrix  # 记录每个男生和每个女生之间的`好感度`
        self._visitedBoy = set()  # 记录每一轮匹配匹配过的男生
        self._visitedGirl = set()  # 记录每一轮匹配匹配过的女生
        self._expBoy = [max(row) for row in adjMatrix]  # 每个男生的期望值
        self._expGirl = [0] * max_  # 每个女生的期望值，为0表示只要有一个男生就可以
        self._slack = []  # 记录每个女生如果能被男生倾心最少还需要多少期望值
        self._row = len(adjMatrix)
        self._col = len(adjMatrix[0])

    def getResult(self) -> int:
        """
        每一轮匹配从左侧男生开始,为每个男生找对象
        每次都从右侧第一个女生开始,选择一个女生,使男女两人的期望和要等于两人之间的好感度。
        每一轮匹配,每个女生只会被尝试匹配一次!
        为每个男生解决归宿问题的方法是:如果找不到就降低期望值,直到找到为止
        """
        for boy in range(self._row):
            self._slack = [INF] * self._col
            # 记录每轮匹配中男生女生是否被尝试匹配过
            while True:
                self._visitedBoy.clear()
                self._visitedGirl.clear()
                # 找到归宿 退出
                if self._find_path(boy):
                    break
                else:
                    # 如果不能找到 就降低期望值
                    # 最小可降低的期望值
                    delta = INF
                    for c in range(self._col):
                        if c not in self._visitedGirl and delta > self._slack[c]:
                            delta = self._slack[c]
                    for r in range(self._row):
                        if r in self._visitedBoy:
                            # 所有访问过的男生降低期望值
                            self._expBoy[r] -= delta
                    for c in range(self._col):
                        if c in self._visitedGirl:
                            # 所有访问过的女生增加期望值
                            self._expGirl[c] += delta
                        else:
                            self._slack[c] -= delta

        # 匹配完成 求出所有配对的好感度的和
        res = 0
        for girl, boy in enumerate(self._match):
            if boy != -1:
                res += self._graph[boy][girl]
        return res

    def _find_path(self, boy: int) -> bool:
        self._visitedBoy.add(boy)
        for girl in range(self._col):
            if girl in self._visitedGirl:
                continue
            tmp_delta = self._expBoy[boy] + self._expGirl[girl] - self._graph[boy][girl]
            # 符合要求
            if tmp_delta == 0:
                self._visitedGirl.add(girl)
                if self._match[girl] == -1 or self._find_path(self._match[girl]):
                    self._match[girl] = boy
                    return True
            # 女生要得到男生的倾心 还需多少期望值
            elif self._slack[girl] > tmp_delta:
                self._slack[girl] = tmp_delta

        return False


if __name__ == '__main__':
    # 2172. 数组的最大与和
    class Solution:
        def maximumANDSum(self, nums: List[int], numSlots: int) -> int:
            slots = list(range(1, numSlots + 1)) + list(range(1, numSlots + 1))
            adjMatrix = [[0 for _ in slots] for _ in nums]
            for i in range(len(nums)):
                for j in range(numSlots * 2):
                    adjMatrix[i][j] = nums[i] & slots[j]
            return KM(adjMatrix).getResult()

    print(Solution().maximumANDSum([1, 2, 3, 4, 5, 6], 3))
    print(Solution().maximumANDSum(nums=[1, 3, 10, 4, 7, 1], numSlots=9))
    print(Solution().maximumANDSum(nums=[14, 7, 9, 8, 2, 4, 11, 1, 9], numSlots=8))
