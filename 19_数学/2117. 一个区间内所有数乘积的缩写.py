# 给你两个正整数 left 和 right ，满足 left <= right 。请你计算 闭区间 [left, right] 中所有整数的 乘积 。
# 由于乘积可能非常大，你需要将它按照以下步骤 缩写 ：
# https://leetcode-cn.com/problems/abbreviating-the-product-of-a-range/


class Solution:
    def abbreviateProduct(self, left: int, right: int) -> str:
        """请你返回一个字符串，表示 闭区间 [left, right] 中所有整数 乘积 的 缩写 。"""
        product = 1
        for num in range(left, right + 1):
            product *= num
        product = str(product)

        n = len(product)
        removed = product.rstrip('0')
        zeros = n - len(removed)

        if len(removed) <= 10:
            return removed + 'e' + str(zeros)
        else:
            return removed[:5] + '...' + removed[-5:] + 'e' + str(zeros)

