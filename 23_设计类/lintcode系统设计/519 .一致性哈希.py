"""
一般的数据库进行horizontal shard的方法是指
把 id 对 数据库服务器总数 n 取模,然后来得到他在哪台机器上。
这种方法的缺点是,当数据继续增加,我们需要增加数据库服务器
将 n 变为 n+1 时,几乎所有的数据都要移动,这就造成了不 consistent。
为了减少这种 naive 的 hash方法(%n) 带来的缺陷
出现了一种新的hash算法:一致性哈希的算法——Consistent Hashing
这种算法有很多种实现方式,这里我们来实现一种简单的 Consistent Hashing。

1.将 id 对 360 取模,假如一开始有3台机器,
那么让3台机器分别负责0~119, 120~239, 240~359 的三个部分。
那么模出来是多少,查一下在哪个区间,就去哪台机器。
2.当机器从 n 台变为 n+1 台了以后,我们从n个区间中,
找到最大的一个区间,然后一分为二,把一半给第n+1台机器。
3.比如从3台变4台的时候,我们找到了第3个区间0~119是当前最大的一个区间,
那么我们把0~119分为0~59和60~119两个部分。
0~59仍然给第1台机器,60~119给第4台机器。
4.然后接着从4台变5台,我们找到最大的区间是第3个区间120~239,
一分为二之后,变为 120~179, 180~239。

假设一开始所有的数据都在一台机器上，请问加到第 n 台机器的时候，
区间的分布情况和对应的机器编号分别是多少？
当最大区间出现多个时，我们拆分编号较小的那台机器。
"""

from typing import List
from sortedcontainers import SortedList


class Solution:
    def __init__(self) -> None:
        # !也可以用堆来维护最大区间
        self._intervals = SortedList([(0, 359, 1)], key=lambda x: (-(x[1] - x[0]), x[2]))

    def consistent_hashing(self, n: int) -> List[List[int]]:
        """
        @param n: a positive integer
        @return: n x 3 matrix
                we will sort your return value in output
        """
        sl = self._intervals
        for i in range(2, n + 1):
            max_ = sl.pop(0)

            start, end, preI = max_
            new1 = (start, (end + start) // 2, preI)
            new2 = ((end + start) // 2 + 1, end, i)
            sl.add(new1)
            sl.add(new2)

        res = []
        for inter in sl:
            res.append(list(inter))
        res.sort(key=lambda x: x[0])
        return res


print(Solution().consistent_hashing(3))
