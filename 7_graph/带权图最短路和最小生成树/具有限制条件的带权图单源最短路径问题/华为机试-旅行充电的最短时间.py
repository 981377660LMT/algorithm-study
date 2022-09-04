#############################################################
# 张三购买了一辆续航里程数达1000公里的某自动驾驶新能源车。
# 某一天车辆充满电后，需从甲城出发前往距离D公里远的乙域，全程走高速。
# 车载导航提示沿途有N个休息站均可提供充电服务，各休息站均可实时提供当前充电排队时间(小时)。
# 请协·助规划时间最优的休息站充电方案，返回最短的旅行用时。
# 为方便计算，高速上的行驶速度固定为100公里/小时。
# 规划时可不必考虑保留安全续航里程数，汽车可以将电完全用光，
# 1000公里续航的汽车按100公里/小时，可以开10个小时。
# 每次充电时间固定为1小时，完成后电量充满。
# 各站点充电排队时间不会变化，充电排队过程不耗电。
# 请协助规划时间最优的休息站充电方案，返回最短的旅行用时。

# 每次充电时间固定为1小时，完成后电量充满。
# !电动车游览城市 每个点转移到下一个点 或者充电
# !注意到是无环图 没必要dijk 可以直接dp 少一个log (不过懒得改了hh)

from collections import defaultdict
from heapq import heappop, heappush

INF = int(1e18)
D = int(input())  # 距离D公里远
N = int(input())  # 沿途有N个休息站
stops = {}
points = [0]
dist = defaultdict(lambda: INF)  # (当前位置,油量)
for _ in range(N):
    pos, wait = map(int, input().split())
    stops[pos] = wait + 1  # 距离和充电时间(小时)
    points.append(pos)
points.append(D)

pq = [(0, 0, 1000)]
while pq:
    cur_cost, cur_id, cur_fuel = heappop(pq)
    if cur_id == len(points) - 1:
        print(int(cur_cost))
        exit(0)

    if cur_cost > dist[(cur_id, cur_fuel)]:
        continue

    # 充电
    if points[cur_id] in stops:
        time = stops[points[cur_id]]
        if time + cur_cost < dist[(cur_id, 1000)]:
            dist[(cur_id, 1000)] = time + cur_cost
            heappush(pq, (time + cur_cost, cur_id, 1000))

    # 走到下个充电站
    weight = points[cur_id + 1] - points[cur_id]
    if cur_fuel - weight >= 0:
        nextFuel = cur_fuel - weight
        if cur_cost + weight / 100 < dist[(cur_id, nextFuel)]:
            dist[(cur_id, nextFuel)] = cur_cost + weight / 100
            heappush(pq, (cur_cost + weight / 100, cur_id + 1, nextFuel))

print(-1)
#########################################
