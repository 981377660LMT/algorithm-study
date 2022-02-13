MOD = 1000000007

#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# @param arr int整型一维数组
# @param a int整型
# @param b int整型
# @return int整型
#
class Solution:
    def countTriplets(self, arr: list[int], a: int, b: int):
        # write code here
        n = len(arr)
        res = 0
        for mid in range(n):
            leftCount = 0
            for left in range(mid):
                if abs(arr[left] - arr[mid]) <= a:
                    leftCount += 1

            rightCount = 0
            for right in range(mid + 1, n):
                if abs(arr[right] - arr[mid]) <= b:
                    rightCount += 1
            res += leftCount * rightCount
            res %= MOD
        return res
