# 3614. 用特殊操作处理字符串 II
# https://leetcode.cn/problems/process-string-with-special-operations-ii/description/
#
# 给你一个字符串 s，由小写英文字母和特殊字符：'*'、'#' 和 '%' 组成。
#
# 同时给你一个整数 k。
#
# 请根据以下规则从左到右处理 s 中每个字符，构造一个新的字符串 result：
#
# 如果字符是 小写 英文字母，则将其添加到 result 中。
# 字符 '*' 会 删除 result 中的最后一个字符（如果存在）。
# 字符 '#' 会 复制 当前的 result 并追加到其自身后面。
# 字符 '%' 会 反转 当前的 result。
# 返回最终字符串 result 中第 k 个字符（下标从 0 开始）。如果 k 超出 result 的下标索引范围，则返回 '.'。
#
# 倒着处理.
#
# 添加、删除、复制、反转


class Solution:
    def processStr(self, s: str, k: int) -> str:
        n = len(s)
        size = [0] * n
        ptr = 0
        for i, c in enumerate(s):
            if c == "*":
                ptr = max(0, ptr - 1)
            elif c == "#":
                ptr *= 2
            elif c.islower():
                ptr += 1
            size[i] = ptr

        if k >= size[-1]:  # 下标越界
            return "."

        for c, size_ in zip(s[::-1], size[::-1]):
            if c == "#":
                if k >= size_ // 2:  # k 在复制后的右半边
                    k -= size_ // 2
            elif c == "%":
                k = size_ - 1 - k
            elif c.islower():
                if k == size_ - 1:
                    return c

        raise ValueError("unreachable")
