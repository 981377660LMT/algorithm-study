n = int(input())
A = [0] + [int(i) for i in input().split()]
B = [0] + [-1] * n
for _ in range(int(input())):
    line = input().split()
    if line[0] == '1':
        k, x, y = int(line[1]), int(line[2]), int(line[3])
        B[y : y + k] = A[x : x + k]
    else:
        print(B[int(line[1])])


# Python 切片为什么不会越界
# 当我们根据单个索引进行取值时，如果索引越界，就会得到报错：“IndexError: list index out of range”。

# 当左索引值大于等于右索引值时，切片结果为空对象。
# 当左或右索引值大于序列的长度值时，就用长度值作为该索引值；
# 当左索引值缺省或者为 None 时，就用 0 作为左索引值；
# 当右索引值缺省或者为 None 时，就用序列长度值作为右索引值；
