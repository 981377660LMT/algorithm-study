#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# 返回在所有合法的三角形的组成中，最大的三角形的周长减去最小的三角形的周长的值
# @param n int整型 代表题目中的n
# @param a int整型一维数组 代表n个数的大小
# @return int整型
#
from typing import List


class Solution:
    def formTriangle(self, n: int, a: List[int]) -> int:
        """
        从n个数中找出三个数来组成一个三角形
        周长最大的三角形的周长减去周长最小的三角形的周长是多少

        排好序后，最大值肯定是连续的3个数
        最小值的话，最长的两条边一定是连续的
        """
        a = sorted(a)
        max_ = 0
        for i in range(n - 1, 1, -1):
            if a[i - 2] + a[i - 1] > a[i]:
                max_ = a[i - 2] + a[i - 1] + a[i]
                break

        min_ = int(1e20)
        ok = False
        for mid in range(1, n - 1):
            for left in range(mid):
                if a[left] + a[mid] > a[mid + 1]:
                    min_ = a[left] + a[mid] + a[mid + 1]
                    ok = True
                    break
            if ok:
                break
        return max_ - min_

