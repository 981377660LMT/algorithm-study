nums = [1, 2, 3, 4]
target = [4, 1, 3, 2]
stack = []
index = 0
for num in nums:
    stack.append(num)
    while stack and index < len(target) and target[index] == stack[-1]:
        stack.pop()
        index += 1

print('Yes' if len(stack) == 0 else 'No')

