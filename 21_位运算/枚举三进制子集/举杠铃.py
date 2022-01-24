# n个杠铃，m个铁片
# 每个铁片可以加在杠铃的左边，或者右边，或者不加
# 但如果加的话，杠铃两边需要保持重量相同
# 然后问，所有合法的杠铃（可不加铁片或加相同重量的铁片）的重量


# 三进制枚举，
# 第i位为1代表给左边加第i块铁片，为2代表给右边加该铁片，为0代表不加
# 这样所有的集合都被枚举到了，
# n,m<=14；plates[i],barbells[i]<=10^8

# 时间复杂度:
# 3^n*n

# 与二进制枚举的区别：
# `state = 1<<m` 改为 `state = pow(m,3)`
# `(state>>j)&1` 改为 `state%3==0/1/2`
from typing import List


class Solution:
    def Barbells(self, n: int, m: int, plates: List[int], barbells: List[int]) -> List[int]:
        res = set()

        for state in range(3 ** m):
            leftSum, rightSum = 0, 0
            # 检查state的第i位是0/1/2中的哪个
            for i in range(m):
                prefix = state // (3 ** i)
                mod = prefix % 3
                if mod == 1:
                    leftSum += plates[i]
                elif mod == 2:
                    rightSum += plates[i]

            if leftSum == rightSum:
                res |= {barbells[i] + leftSum + rightSum for i in range(n)}

        return list(res)


print(Solution().Barbells(3, 3, [1, 2, 3], [4, 5, 6]))

