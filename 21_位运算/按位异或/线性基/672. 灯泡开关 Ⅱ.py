# 房间中有 n 只已经打开的灯泡，
# !编号从 1 到 n 。墙上挂着 4 个开关 。

# 这 4 个开关各自都具有不同的功能，其中：

# 开关 1 ：反转当前所有灯的状态（即开变为关，关变为开）
# 开关 2 ：反转编号为偶数的灯的状态（即 2, 4, ...）
# 开关 3 ：反转编号为奇数的灯的状态（即 1, 3, ...）
# 开关 4 ：反转编号为 j = 3k + 1 的灯的状态，其中 k = 0, 1, 2, ...（即 1, 4, 7, 10, ...）
# 你必须 恰好 按压开关 presses 次。每次按压，你都需要从 4 个开关中选出一个来执行按压操作。

# !给你两个整数 n 和 presses ，执行完所有按压之后，返回 不同可能状态 的数量。


from itertools import product


class Solution:
    def flipLights(self, n: int, presses: int) -> int:
        """
        只有四个开关,不需要求线性基什么的
        前6个灯唯一地决定了其余的灯。
        这是因为每一个修改 第 x 的灯光的操作都会修改第 (x+6) 的灯光。
        注意到灯泡的状态六个为一组，直接枚举即可
        """
        count = min(n, 6)
        visited = set()
        for cand in product([0, 1], repeat=4):
            sum_ = sum(cand)
            if (sum_ & 1 == presses & 1) and sum_ <= presses:
                status = [1] * count
                for i in range(count):
                    if cand[0]:
                        status[i] ^= 1
                    if cand[1] and i % 2 == 0:
                        status[i] ^= 1
                    if cand[2] and i % 2 == 1:
                        status[i] ^= 1
                    if cand[3] and i % 3 == 0:
                        status[i] ^= 1
                visited.add(tuple(status))

        return len(visited)


print(Solution().flipLights(3, 1))
