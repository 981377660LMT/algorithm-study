"""
在 Consistent Hashing I 中我们介绍了一个比较简单的一致性哈希算法,这个简单的版本有两个缺陷：
1.增加一台机器之后,数据全部从其中一台机器过来,这一台机器的读负载过大,对正常的服务会造成影响。
2.当增加到3台机器的时候,每台服务器的负载量不均衡,为1:1:2。

为了解决这个问题，引入了 micro-shards 的概念，一个更好的算法是这样：
1.将 360° 的区间分得更细。从 0~359 变为一个 0 ~ n-1 的区间，将这个区间首尾相接，连成一个环。
2.当加入一台新的机器的时候，随机选择在环上撒 k 个点，代表这台机器的 k 个 micro-shards。
3.每个数据在环上也对应一个点，这个点通过一个 hash function 来计算。
4.一个数据该属于哪台机器负责管理，
是按照该数据对应的环上的点在环上顺时针碰到的第一个 micro-shard 点所属的机器来决定。

n 和 k在真实的 NoSQL 数据库中一般是 2^64 和 1000。
"""


from random import randint
from typing import List
from sortedcontainers import SortedList


class Solution:
    @classmethod
    def create(cls, n: int, k: int) -> "Solution":
        """facade"""
        return cls(n, k)

    def __init__(self, cycleLength: int, microShardCount: int) -> None:
        """n 和 k在真实的 NoSQL 数据库中一般是 2^64 和 1000。"""
        self.cycleLength = cycleLength
        self.microShardCount = microShardCount
        self.machinePosition = SortedList(key=lambda x: x[0])
        self._visited = set()

    def addMachine(self, machine_id: int) -> List[int]:
        """
        添加新机器,返回碎片ID列表

        当加入一台新的机器的时候，随机选择在环上撒 k 个点，代表这台机器的 k 个 micro-shards
        """

        n, k, cur = self.cycleLength, self.microShardCount, []
        while len(cur) < k:
            rand = randint(0, n - 1)
            if rand not in self._visited:
                self._visited.add(rand)
                cur.append(rand)
                self.machinePosition.add((rand, machine_id))
        return cur

    def getMachineIdByHashCode(self, hashcode: int) -> int:
        """
        返回机器id

        由环上顺时针碰到的第一个 micro-shard 点所属的机器来决定
        """
        pos = self.machinePosition.bisect_left((hashcode,))
        if pos == len(self.machinePosition):
            pos = 0
        return self.machinePosition[pos][1]


if __name__ == "__main__":
    S = Solution.create(100, 3)
    print(S.addMachine(1))  # [77,83,86]
    print(S.getMachineIdByHashCode(4))  # 1
    print(S.addMachine(2))  # [15,35,93]
    print(S.getMachineIdByHashCode(61))  # 1
    print(S.getMachineIdByHashCode(91))  # 2
