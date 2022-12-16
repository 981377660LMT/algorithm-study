# B - FizzBuzz Sum-等差数列求和
# 一个数列,其中第i项(i>=1)为:
# i%3==0 and i%5==0: FizzBuzz
# i%3==0: Fizz
# i%5==0: Buzz
# i: i
# !求出前n项中包含的数字的和 (1<=n<=1e9)

# !使用周期+等差数列求和加速模拟

preSum = [0]  # 第一个周期的前缀和
digitCount = 0
for i in range(15):
    if i % 3 == 0 or i % 5 == 0:
        preSum.append(preSum[-1])
    else:
        preSum.append(preSum[-1] + i)
        digitCount += 1


def arithmeticSum2(first: int, diff: int, item: int) -> int:
    """等差数列求和 first:首项 diff:公差 item:项数"""
    last = first + (item - 1) * diff
    return item * (first + last) // 2


def fizzBuzzSum(n: int) -> int:
    div, mod = divmod(n, 15)
    sum1 = arithmeticSum2(preSum[-1], 15 * digitCount, div)
    sum2 = 0
    for i in range(1, mod + 1):
        if i % 3 == 0 or i % 5 == 0:
            continue
        sum2 += i + 15 * div
    return sum1 + sum2


n = int(input())
print(fizzBuzzSum(n))
