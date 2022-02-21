# 在[0,n][0,n]范围中，选取一个最大的数x，满足x% a = b，不过这个范围可能会很大。
# 给定如上所述的a , b ,n，返回满足条件的最大的x。
#
class Solution:
    def solve(self, a, b, n):
        # write code here
        return (n - b) // a * a + b

