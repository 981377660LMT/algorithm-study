# 每次操作会使排列变成它的前一个排列，
# 题目等价于求从小到大数第几个排列-1，
# 或者相当于求有多少个排列比它小
# !求在所有的`组合中`，当前这个是字典序第k小 (k从0开始)
# 有点类似康托展开的思想
# 1 <= s.length <= 3000

from 康托展开_取模 import calRank


class Solution:
    def makeStringSorted(self, s: str) -> int:
        """求在所有`排列`中,当前这个是字典序第几小(rank>=0)"""
        return calRank(s)


if __name__ == "__main__":

    assert Solution().makeStringSorted(s="abc") == 0
    assert Solution().makeStringSorted(s="cba") == 5
    assert Solution().makeStringSorted(s="ccba") == 11
    # 输出：5
    # 解释：模拟过程如下所示：
    # 操作 1：i=2，j=2。交换 s[1] 和 s[2] 得到 s="cab" ，然后反转下标从 2 开始的后缀字符串，得到 s="cab" 。
    # 操作 2：i=1，j=2。交换 s[0] 和 s[2] 得到 s="bac" ，然后反转下标从 1 开始的后缀字符串，得到 s="bca" 。
    # 操作 3：i=2，j=2。交换 s[1] 和 s[2] 得到 s="bac" ，然后反转下标从 2 开始的后缀字符串，得到 s="bac" 。
    # 操作 4：i=1，j=1。交换 s[0] 和 s[1] 得到 s="abc" ，然后反转下标从 1 开始的后缀字符串，得到 s="acb" 。
    # 操作 5：i=2，j=2。交换 s[1] 和 s[2] 得到 s="abc" ，然后反转下标从 2 开始的后缀字符串，得到 s="abc" 。
