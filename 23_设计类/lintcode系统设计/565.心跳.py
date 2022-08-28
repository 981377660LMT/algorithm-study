# 在Master-Slave结构模式中，
# slave端会每隔k秒向master端发送ping请求表示自己还在运行。
# 如果master端在 2 * k 秒内没有接收到任何来自slave的ping请求，
# 那么master端会向管理员发送一个警告(比如发送一个电子邮件)。

from typing import List


class HeartBeat:
    def __init__(self):
        self.preTime = dict()  # slave 上一次上报的时间

    def initialize(self, slaves_ip_list: List[str], k: int) -> None:
        self.interval = k
        for ip in slaves_ip_list:
            self.preTime[ip] = 0

    def ping(self, timestamp: int, slave_ip: str) -> None:
        """master 端从 slave 端收到 ping 请求"""
        if slave_ip not in self.preTime:
            return
        self.preTime[slave_ip] = timestamp

    def getDiedSlaves(self, timestamp: int) -> List[str]:
        """这个方法会定期的(两次执行之间的时间间隔不确定)执行"""
        res = []
        for slaveId, pre in self.preTime.items():
            if timestamp - pre >= 2 * self.interval:
                res.append(slaveId)
        return res


if __name__ == "__main__":
    hb = HeartBeat()
    hb.initialize(["10.173.0.2", "10.173.0.3"], 10)
    hb.ping(1, "10.173.0.2")
    hb.getDiedSlaves(20)
    hb.getDiedSlaves(21)
    hb.ping(22, "10.173.0.2")
    hb.ping(23, "10.173.0.3")
    hb.getDiedSlaves(24)
    hb.getDiedSlaves(42)
