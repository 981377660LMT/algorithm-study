import enum


class Solution:
    def solve(self, s, k):
        """是否存在一个座位使得距离每个人至少k"""
        groups = [len(g) for g in s.split('x')]
        for i, gSize in enumerate(groups):
            if i == 0 or i == len(groups) - 1:
                if gSize >= k:
                    return True
            elif gSize >= 2 * k - 1:
                return True
        return False

