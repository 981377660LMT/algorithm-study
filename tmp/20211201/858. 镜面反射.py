# 正方形房间的墙壁长度为 p，一束激光从西南角射出，首先会与东墙相遇，入射点到接收器 0 的距离为 q 。
# 返回光线最先遇到的接收器的编号（保证光线最终会遇到一个接收器）。

# 思路：延展镜面，最终相当于从一个举行的左下角走到右上角
class Solution:
    def mirrorReflection(self, p: int, q: int) -> int:
        k = 1
        while q * k % p != 0:
            k += 1

        times = q * k / p % 2  # 上下延伸的段数奇偶
        dir = k % 2  # 左右

        if times == 1 and dir == 1:
            return 1
        if times == 0 and dir == 1:
            return 0
        if times == 1 and dir == 0:
            return 2

        return -1


print(Solution().mirrorReflection(p=2, q=1))
