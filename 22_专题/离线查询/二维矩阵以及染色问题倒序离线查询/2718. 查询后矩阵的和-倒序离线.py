# 给你一个整数 n 和一个下标从 0 开始的 二维数组 queries ，其中 queries[i] = [typei, indexi, vali] 。

# 一开始，给你一个下标从 0 开始的 n x n 矩阵，所有元素均为 0 。每一个查询，你需要执行以下操作之一：


# 如果 typei == 0 ，将第 indexi 行的元素全部修改为 vali ，覆盖任何之前的值。
# 如果 typei == 1 ，将第 indexi 列的元素全部修改为 vali ，覆盖任何之前的值。
# 请你执行完所有查询以后，返回矩阵中所有整数的和。


# 在线的解法：https://leetcode.cn/problems/sum-of-matrix-after-queries/solution/shu-zhuang-shu-zu-shi-jian-chuo-zai-xian-eqgo/

from typing import List


class Solution:
    def matrixSumQueries(self, n: int, queries: List[List[int]]) -> int:
        res = 0
        visited = [set(), set()]
        for type_, index, value in queries[::-1]:
            if index not in visited[type_]:
                visited[type_].add(index)
                remain = n - len(visited[type_ ^ 1])
                res += remain * value
        return res
