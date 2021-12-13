from typing import List

# 1 <= customers.length <= 10^5

# arrivali 是第 i 位顾客到达的时间，到达时间按 非递减 顺序排列。
# timei 是给第 i 位顾客做菜需要的时间。
class Solution:
    def averageWaitingTime(self, customers: List[List[int]]) -> float:
        waitSum, finishTime = 0, 0
        for start, cost in customers:
            finishTime = max(finishTime, start) + cost
            waitSum += finishTime - start
        return waitSum / len(customers)


print(Solution().averageWaitingTime(customers=[[5, 2], [5, 4], [10, 3], [20, 1]]))
# 输出：3.25000
# 解释：
# 1) 第一位顾客在时刻 5 到达，厨师拿到他的订单并在时刻 5 立马开始做菜，并在时刻 7 完成，第一位顾客等待时间为 7 - 5 = 2 。
# 2) 第二位顾客在时刻 5 到达，厨师在时刻 7 开始为他做菜，并在时刻 11 完成，第二位顾客等待时间为 11 - 5 = 6 。
# 3) 第三位顾客在时刻 10 到达，厨师在时刻 11 开始为他做菜，并在时刻 14 完成，第三位顾客等待时间为 14 - 10 = 4 。
# 4) 第四位顾客在时刻 20 到达，厨师拿到他的订单并在时刻 20 立马开始做菜，并在时刻 21 完成，第四位顾客等待时间为 21 - 20 = 1 。
# 平均等待时间为 (2 + 6 + 4 + 1) / 4 = 3.25 。
