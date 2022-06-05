# 编辑器共有五种指令，如下：

# 1、I x，在光标处插入数值 x。
# 2、D，将光标前面的第一个元素删除，如果前面没有元素，则忽略此操作。
# 3、L，将光标向左移动，跳过一个元素，如果左边没有元素，则忽略此操作。
# 4、R，将光标向右移动，跳过一个元素，如果右边没有元素，则忽略此操作。
# 5、Q k，假设此刻光标之前的序列为 a1,a2,…,an，输出 max1≤i≤kSi，其中 Si=a1+a2+…+ai。


# 对顶栈
# 【。。。。   。。。。。】

import sys

input = sys.stdin.readline


q = int(input())

stack1, stack2 = [], []
preSum = [0] * q
preSumMax = [-int(1e20)] * q

for _ in range(q):
    opt, *rest = input().split()
    if opt == 'I':
        num = int(rest[0])
        stack1.append(num)
        preSum[len(stack1)] = preSum[len(stack1) - 1] + num
        preSumMax[len(stack1)] = max(preSumMax[len(stack1) - 1], preSum[len(stack1)])
    elif opt == 'D':
        if stack1:
            stack1.pop()
    elif opt == 'L':
        if stack1:
            stack2.append(stack1.pop())
    elif opt == 'R':
        if stack2:
            top = stack2.pop()
            stack1.append(top)
            preSum[len(stack1)] = preSum[len(stack1) - 1] + top
            preSumMax[len(stack1)] = max(preSumMax[len(stack1) - 1], preSum[len(stack1)])
    else:
        k = int(rest[0])
        print(preSumMax[k])

