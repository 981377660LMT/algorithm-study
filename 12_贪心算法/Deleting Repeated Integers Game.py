# 最快删除repeat number

# 奇数长度的数组
# repeat数字个数 len(A) - len(set(A))
class Solution:
    def solve(self, A):
        repeats = len(A) - len(set(A))
        return repeats + 1 >> 1

