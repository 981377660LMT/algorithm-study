#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# ​返回按照这些花排成一个圆的序列中最小的“丑陋度”
# @param n int整型 花的数量
# @param array int整型一维数组 花的高度数组
# @return int整型
# n朵花排成一圈，最小化相邻两朵花高度差的最大值，输出最大值。
# 最后肯定是山脉数组
#
from typing import List

# 最优策略为：第1小的数左侧依次为第2、4、6···小的数，右侧依次为3、5、7···小的数。
# 需要反证得出结果
class Solution:
    def arrangeFlowers(self, n: int, array: List[int]):
        if n <= 3:
            return max(array) - min(array)
        array.sort()
        return max(array[i + 2] - array[i] for i in range(n - 2))
