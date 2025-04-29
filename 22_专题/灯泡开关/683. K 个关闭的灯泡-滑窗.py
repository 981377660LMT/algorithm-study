# 683. K 个关闭的灯泡
# https://leetcode.cn/problems/k-empty-slots/description/
#
# n 个灯泡排成一行，编号从 1 到 n 。最初，所有灯泡都关闭。
# 每天 只打开一个 灯泡，直到 n 天后所有灯泡都打开。
# 给你一个长度为 n 的灯泡数组 blubs ，其中 bulbs[i] = x 意味着在第 (i+1) 天，我们会把在位置 x 的灯泡打开，其中 i 从 0 开始，x 从 1 开始。
# 给你一个整数 k ，请返回恰好有两个打开的灯泡，且它们中间 正好 有 k 个 全部关闭的 灯泡的 最小的天数 。
# 如果不存在这种情况，返回 -1 。


from typing import List

INF = int(1e18)


class Solution:
    def kEmptySlots(self, bulbs: List[int], k: int) -> int:
        n = len(bulbs)
        days = [0] * n  # 每个灯泡在第几天被打开
        for day, pos in enumerate(bulbs, start=1):
            days[pos - 1] = day

        res = INF
        left, right = 0, k + 1
        while right < n:
            # 检查 (left, right) 中间的所有位置 j
            # 如果某 j 对应 days[j] <= max(days[left], days[right])
            # 则说明 j 那盏灯先于两端之一被打开，不满足“中间恰好 k 盏关着”的要求
            for i in range(left + 1, right):
                if days[i] < days[left] or days[i] < days[right]:
                    # 直接把 i 跳到 j，跳过那些不可能的左端
                    left, right = i, i + +k + 1
                    break
            else:
                # 找到一个满足条件的窗口，[left, right] 两端打开，中间 k 盏一直关着
                res = min(res, max(days[left], days[right]))
                left, right = right, right + k + 1

        return -1 if res == INF else res


if __name__ == "__main__":
    sol = Solution()
    print(sol.kEmptySlots([1, 3, 2], 1))  # 输出 2：第2天开3号，再有1盏(2号)关着
    print(sol.kEmptySlots([1, 2, 3], 1))  # 输出 -1：始终没有恰好中间1盏关着
    print(
        sol.kEmptySlots([3, 1, 5, 4, 2], 2)
    )  # 输出 3：当天序 [3,1,5,4,2] 中第3天开5号，与3号中间2盏(4、2)关着
