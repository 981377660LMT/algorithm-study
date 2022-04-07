# 找到下一个最近的合题意的时间
class Solution:
    def solve(self, s: str):
        def check(time: int):
            h = str(time // 60).zfill(2)
            m = str(time % 60).zfill(2)
            return all(c in visited for c in h + m)

        visited = set(s)

        h, m = map(int, s.split(":"))
        curTime = 60 * h + m
        nextTime = (curTime + 1) % 1440

        # 每次加1暴力搜索
        while not check(nextTime):
            nextTime = (nextTime + 1) % 1440

        nextH = str(nextTime // 60).zfill(2)
        nextM = str(nextTime % 60).zfill(2)
        return nextH + ":" + nextM
