# !1279. 红绿灯路口
# https://leetcode.cn/problems/traffic-light-controlled-intersection/description/?envType=problem-list-v2&envId=concurrency
#
# 实现一个红绿灯控制系统，确保不会发生死锁。
# 两条路上的红绿灯不可以同时为绿灯。这意味着，当 A 路上的绿灯亮起时，B 路上的红灯会亮起；当 B 路上的绿灯亮起时，A 路上的红灯会亮起.
# 开始时，A 路上的绿灯亮起，B 路上的红灯亮起。当一条路上的绿灯亮起时，所有车辆都可以从任意两个方向通过路口，直到另一条路上的绿灯亮起。不同路上的车辆不可以同时通过路口。
#
# 关键点是两条路的红绿灯不能同时为绿灯，而且需要在必要时切换信号灯。
#
# 1. 使用互斥锁确保同一时间只有一个车辆能通过执行决策逻辑
# 2. 记录当前哪条路的绿灯亮起
# 3. 当车辆到达时，检查它所在的道路是否有绿灯：
#    如果有，直接通过
#    如果没有，切换红绿灯，然后通过

import threading
from typing import Callable


class TrafficLight:
    def __init__(self):
        self.greenRoadId = 1  # 1 表示 Road A，2 表示 Road B
        self.lock = threading.Lock()

    def carArrived(
        self,
        carId: int,  # ID of the car
        roadId: int,  # ID of the road the car travels on. Can be 1 (road A) or 2 (road B)
        direction: int,  # Direction of the car
        turnGreen: "Callable[[], None]",  # Use turnGreen() to turn light to green on current road
        crossCar: "Callable[[], None]",  # Use crossCar() to make car cross the intersection
    ) -> None:
        with self.lock:
            if self.greenRoadId != roadId:
                turnGreen()
                self.greenRoadId = roadId
            crossCar()
