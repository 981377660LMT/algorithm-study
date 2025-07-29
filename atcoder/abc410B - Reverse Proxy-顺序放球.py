#
# !顺序放球
#
# 有 N 个编号为 1…N 的箱子，开始时都为空。
# 接下来依次到来 Q 个小球；第 i 个小球给定指令 Xi：
#
# • Xi ≥ 1：把该球放入编号为 Xi 的箱子。
# • Xi = 0：把该球放入当前球数最少的箱子；若有多个并列，选编号最小的那一个。
#
# 请输出每个小球最终被放入的箱子编号。
#
# 桶加指针的懒删除保证每只球和每个桶指针最多前移一次，总体 O(N+Q)。
if __name__ == "__main__":

    def solve():
        N, _ = map(int, input().split())
        X = list(map(int, input().split()))

        counter = [0] * (N + 1)
        ptr, minCount = 1, 1
        for x in X:
            if x >= 1:
                counter[x] += 1
                print(x)
            else:
                while counter[ptr] >= minCount:
                    ptr += 1
                    if ptr == N + 1:
                        ptr = 1
                        minCount += 1
                counter[ptr] += 1
                print(ptr)

    solve()
