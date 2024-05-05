def kthSuperset(mask: int, k: int) -> int:
    """
    mask的第k个超集,k 从0开始.
    即将k填入mask的二进制中为0的部分, bmi2指令集中有pdep指令可以完成这个操作.
    """
    i, j = 0, 0
    while k >> j:
        if (mask >> i & 1) == 0:
            mask |= (k >> j & 1) << i
            j += 1
        i += 1
    return mask


class Solution:
    # 3133. 数组最后一个元素的最小值
    # https://leetcode.cn/problems/minimum-array-end/description/
    # !构造一个长度为n，按位与为x的严格递增数组，返回最后一个元素的最小值.
    #
    # 给你两个整数 n 和 x 。你需要构造一个长度为 n 的 正整数 数组 nums ，
    # 对于所有 0 <= i < n - 1 ，满足 nums[i + 1] 大于 nums[i] ，
    # !并且数组 nums 中所有元素的按位 AND 运算结果为 x 。
    # 返回 nums[n - 1] 可能的 最小 值。

    # !从集合的视角看，and_ 是每个 nums[i] 的子集。换句话说，nums[i] 一定是 and_ 的超集。
    # 为了让 nums[i] 尽可能小，我们应该选择and_的超集中最小的n个数.
    # !那么就变成了找到第k小的超集.
    #
    # 将 n-1 填进 and_  的二进制中为 0 的部分即可
    # bmi2 指令集中刚好有 pdep 这个指令可以完成这个操作
    # https://leetcode.cn/problems/minimum-array-end/solutions/2759297/o1jie-fa-li-yong-bmizhi-ling-by-vclip-pe09/
    def minEnd(self, n: int, and_: int) -> int:
        return kthSuperset(and_, n - 1)
