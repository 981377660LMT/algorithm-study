# 如果n是奇数 最开始的根要放奇数
# 如果父亲是奇数 儿子必须放偶数
# 如果父亲是偶数 儿子尽量放奇数

# 因为parent总是在数组前半部分 所以数组前半部分全放偶数，后半部分全放奇数

n = int(input())
evens = [i for i in range(1, n + 1) if i % 2 == 0]
odds = [i for i in range(1, n + 1) if i % 2 == 1]


for num in evens + odds:
    print(num, end=' ')
