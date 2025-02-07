# 3433. 统计用户被提及情况
# 有两种类型的事件：
# 1. 消息事件：单个用户或者所有被提及，形如 ["MESSAGE", "timestampi", "mentions_stringi"]；
# 2. 离线事件：某个用户离线，在60个单位事件后再次上线，形如 ["OFFLINE", "timestampi", "useri"]。
# 时间戳相同时，离线事件优先于消息事件。
# 统计每个用户被提及的次数。
#
# !用一个数组 nextOnline 标记用户下次在线的时间戳。如果 nextOnline[i]≤ 当前时间戳，则表示用户 i 已在线

from typing import List


class Solution:
    def countMentions(self, numberOfUsers: int, events: List[List[str]]) -> List[int]:
        # 按照时间戳从小到大排序，时间戳相同的，离线事件排在前面
        events.sort(key=lambda e: (int(e[1]), e[0][2]))

        res = [0] * numberOfUsers
        nextOnline = [0] * numberOfUsers
        for kind, time, mention in events:
            time = int(time)
            if kind == "OFFLINE":
                nextOnline[int(mention)] = time + 60
            elif mention == "ALL":
                for i in range(numberOfUsers):
                    res[i] += 1
            elif mention == "HERE":
                for i, t in enumerate(nextOnline):
                    if t <= time:  # 在线
                        res[i] += 1
            else:
                for s in mention.split():
                    id = int(s[2:])
                    res[id] += 1

        return res
