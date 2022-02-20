#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# 返回最少多少次操作后能使这几个数变成一个等差数列，如果完全不能构造成功，就返回-1
# @param n int整型 代表一共有n个数字
# @param b int整型一维数组 代表这n个数字的大小
# @return int整型
#
class Solution:
    def arithmeticSequence(self, n: int, b: int):
        # 使用最少的操作次数，将这几个数构造成一个等差数列。
        """"
        枚举序列 b 中前 2 个数的操作，每个数有 3 种操作（+1，-1，+0），所以共有 9 种情况。
        每种情况确定了一个等差数列，因为等差数列的首元素 a1 和公差 d 已确定。判断后面的数是否符合这个等差数列。
        """
        ...
