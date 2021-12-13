from typing import List

# 请你找出在对 num 的任一`子字符串`执行突变操作（也可以不执行）后，可能得到的 最大整数 ，并用字符串表示返回。
# 子字符串 是字符串中的一个连续序列。
class Solution:
    def maximumNumber(self, num: str, change: List[int]) -> str:
        n = len(num)
        arr = list(num)
        for i in range(n):
            if change[int(arr[i])] > int(arr[i]):
                #  更新右边界
                while i < n and change[int(arr[i])] >= int(arr[i]):
                    arr[i] = str(change[int(arr[i])])
                    i += 1
                break

        return ''.join(arr)


print(Solution().maximumNumber(num="021", change=[9, 4, 3, 5, 7, 2, 1, 9, 0, 6]))
# 输出："934"
# 解释：替换子字符串 "021"：
# - 0 映射为 change[0] = 9 。
# - 2 映射为 change[2] = 3 。
# - 1 映射为 change[1] = 4 。
# 因此，"021" 变为 "934" 。
# "934" 是可以构造的最大整数，所以返回它的字符串表示。
