# https://leetcode.cn/problems/maximum-elegance-of-a-k-length-subsequence/
# 寿司最大美味度
# 给你一个长度为 n 的二维整数数组 items 和一个整数 k 。
# !items[i] = [profiti, categoryi]，
# !其中 profiti 和 categoryi 分别表示第 i 个项目的利润和类别。
# !现定义 items 的 子序列 的 优雅度 可以用
# total_profit + distinct_categories^2 计算，其中 total_profit 是子序列中所有项目的利润总和，distinct_categories 是所选子序列所含的所有类别中不同类别的数量。
# 你的任务是从 items 所有长度为 k 的子序列中，找出 最大优雅度 。
# 用整数形式表示并返回 items 中所有长度恰好为 k 的子序列的最大优雅度。
# 注意：数组的子序列是经由原数组删除一些元素（可能不删除）而产生的新数组，且删除不改变其余元素相对顺序。


# 双参量:一个参数排序，维护另一个量
# !枚举种类数


from typing import List


class Solution:
    def findMaximumElegance(self, items: List[List[int]], k: int) -> int:
        items.sort(key=lambda x: x[0], reverse=True)

        visitedType = set()
        toRemove = []  # 可以删除的数，末尾最小
        curSum = 0
        for num, type_ in items[:k]:
            curSum += num
            if type_ in visitedType:
                toRemove.append(num)
            else:
                visitedType.add(type_)

        count = len(visitedType)
        res = curSum + count * count
        for num, type_ in items[k:]:  # 增加新种类，移除之前最小的数(进来一个新的，出去一个最没用的)
            if not toRemove:
                break
            if type_ in visitedType:
                continue
            curSum += num - toRemove.pop()
            count += 1
            res = max(res, curSum + count * count)
            visitedType.add(type_)
        return res


# items = [[3,2],[5,1],[10,1]], k = 2
print(Solution().findMaximumElegance([[3, 2], [5, 1], [10, 1]], 2))
