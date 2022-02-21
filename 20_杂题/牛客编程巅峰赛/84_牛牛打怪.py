#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# @param n int整型
# @param DEF int整型一维数组
# @return int整型
#
class Solution:
    def Minimumdays(self, n, DEF):
        """
      第i次操作可以选择一个值小于i的元素并使之变为0，
      若没有小于i的元素可以认为一次空操作，次数也是增加
      求使数组元素全为0最少的操作数。
      """
        DEF.sort()
        res = -int(1e20)
        for i in range(n):
            res = max(res + 1, DEF[i])
        return res

