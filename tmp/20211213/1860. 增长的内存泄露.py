from typing import List

# 请你返回一个数组，包含 [crashTime, memory1crash, memory2crash] ，其中 crashTime是程序意外退出的时间（单位为秒）， memory1crash 和 memory2crash 分别是两个内存条最后剩余内存的位数。


class Solution:
    def memLeak(self, memory1: int, memory2: int) -> List[int]:
        time = 1
        while time <= max(memory1, memory2):
            if memory1 < memory2:
                memory2 -= time
            else:
                memory1 -= time
            time += 1
        return [time, memory1, memory2]


print(Solution().memLeak(memory1=8, memory2=11))
