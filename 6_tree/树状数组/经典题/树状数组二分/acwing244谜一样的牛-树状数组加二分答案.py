from BIT import BIT1

# 有 n 头奶牛，已知它们的身高为 1∼n 且各不相同，但不知道每头奶牛的具体身高。
# 现在这 n 头奶牛站成一列，已知第 i 头牛前面有 Ai 头牛比它低，求每头奶牛的身高。
# 输出包含 n 行，每行输出一个整数表示牛的身高。
# 第 i 行输出第 i 头牛的身高。
# n<=10^5

n = int(input())
nums = [0]
for _ in range(n - 1):
    nums.append(int(input()))

res = [1] * n
bit = BIT1(n + 10)
for i in range(1, n + 1):
    bit.add(i, 1)
for i in range(n - 1, -1, -1):
    # 二分答案，如果前面有k个比他矮，则身高是k+1
    # 找到最小的x，满足sum(x)==k+1
    moreThan = nums[i]
    left, right = 1, n
    while left <= right:
        mid = (left + right) >> 1
        if bit.query(mid) >= moreThan + 1:
            right = mid - 1
        else:
            left = mid + 1
    res[i] = left
    bit.add(left, -1)


for num in res:
    print(num)
