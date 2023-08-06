from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的二维整数数组 items 和一个整数 k 。

# items[i] = [profiti, categoryi]，其中 profiti 和 categoryi 分别表示第 i 个项目的利润和类别。

# 现定义 items 的 子序列 的 优雅度 可以用 total_profit + distinct_categories2 计算，其中 total_profit 是子序列中所有项目的利润总和，distinct_categories 是所选子序列所含的所有类别中不同类别的数量。

# 你的任务是从 items 所有长度为 k 的子序列中，找出 最大优雅度 。

# 用整数形式表示并返回 items 中所有长度恰好为 k 的子序列的最大优雅度。


# 注意：数组的子序列是经由原数组删除一些元素（可能不删除）而产生的新数组，且删除不改变其余元素相对顺序。


# 子集
# !枚举种类数

# 可以逆排序后从第k位置双指针，l向左走，r向右走。维护一个Map来记录每种类别的已选取数量，对于每一个r,如果是个新类型，则l向左寻找第一个同类别数量>1的位置，选取r，丢弃l，更新当前结果（当前结果有可能变小，这是允许的，因为后面由于类别数的增加，有可能会变的更大）。过程中不断用当前结果向大更新最终答案。
# 双参数 => 固定一个偏序，维护另一个量


# 作者：YF_Cui
# 链接：https://leetcode.cn/circle/discuss/chtVBq/view/Hy40Ye/
# 来源：力扣（LeetCode）
# 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
class Solution:
    def findMaximumElegance(self, items: List[List[int]], k: int) -> int:
        items.sort(key=lambda x: x[0], reverse=True)
        kind = set()
        nums = []  # 可以删除的数，末尾最小
        curSum = 0
        res = 0
        for s, t in items[:k]:
            curSum += s
            if t in kind:
                nums.append(s)
            else:
                kind.add(t)

        count = len(kind)
        res = curSum + count * count
        for s, t in items[k:]:
            if nums:
                if t in kind:
                    continue
                curSum += s - nums.pop()
                count += 1
                res = max(res, curSum + count * count)
                kind.add(t)
        return res


# items = [[3,2],[5,1],[10,1]], k = 2
print(Solution().findMaximumElegance([[3, 2], [5, 1], [10, 1]], 2))
