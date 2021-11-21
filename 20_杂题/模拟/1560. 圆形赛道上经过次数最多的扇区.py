# 请你以数组形式返回经过次数最多的那几个扇区，按扇区编号 升序 排列。
# 注意，赛道按扇区编号升序逆时针形成一个圆（请参见第一个示例）。

# 注意:中间扇区不用管，都一样
# 只用关心开头结尾
from typing import List


class Solution:
    def mostVisited(self, n: int, rounds: List[int]) -> List[int]:
        s, e = rounds[0], rounds[-1]
        if s <= e:
            # [起点, 终点]
            return list(range(s, e + 1))
        else:
            # [1, 终点]+[起点, n]
            return list(range(1, e + 1)) + list(range(s, n + 1))


print(Solution().mostVisited(4, [1, 3, 1, 2]))
# 输出：[1,2]
# 解释：本场马拉松比赛从扇区 1 开始。经过各个扇区的次序如下所示：
# 1 --> 2 --> 3（阶段 1 结束）--> 4 --> 1（阶段 2 结束）--> 2（阶段 3 结束，即本场马拉松结束）
# 其中，扇区 1 和 2 都经过了两次，它们是经过次数最多的两个扇区。扇区 3 和 4 都只经过了一次。
