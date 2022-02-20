# 有 N 头牛站成一行，被编队为 1、2、3…N，每头牛的身高都为整数。
# 当且仅当两头牛中间的牛身高都比它们矮时，两头牛方可看到对方。
# 现在，我们只知道其中最高的牛是第 P 头，它的身高是 H ，剩余牛的身高未知。
# 但是，我们还知道这群牛之中存在着 M 对关系，每对关系都指明了某两头牛 A 和 B 可以相互看见。
# 求每头牛的身高的最大可能值是多少。

# 一共输出 N 行数据，每行输出一个整数。
# 第 i 行输出的整数代表第 i 头牛可能的最大身高。


n, maxIndex, maxValue, m = map(int, input().split())
diff = [0] * n
canSee = set()
for _ in range(m):
    a, b = sorted(map(int, input().split()))
    if (a, b) in canSee:
        continue
    canSee.add((a, b))
    a, b = a - 1, b - 1
    # a+1到b-1比两端小
    diff[a + 1] -= 1
    diff[b] += 1
for i in range(1, n):
    diff[i] += diff[i - 1]
for less in diff:
    print(less + maxValue)
