# 第 i 天所生成的数组，是由你对第 i-1 天的数组进行如下操作所得的：

# 假如一个元素小于它的左右邻居，那么该元素自增 1。
# 假如一个元素大于它的左右邻居，那么该元素自减 1。
# 首、尾元素 永不 改变。
# 过些时日，你会发现数组将会不再发生变化，请返回最终所得到的数组。

#

from typing import List

# 利用标志位结束循环
class Solution:
    def transformArray(self, arr: List[int]) -> List[int]:
        n = len(arr)
        while True:
            changed = False
            copy = arr[:]
            for i in range(1, n - 1):
                if copy[i - 1] < copy[i] > copy[i + 1]:
                    arr[i] -= 1
                    changed = True
                elif copy[i - 1] > copy[i] < copy[i + 1]:
                    arr[i] += 1
                    changed = True
            if not changed:
                break
        return arr

