from array import array


arr = array("I", bytes(8))

print(arr[0])  # 0
print(arr[1])  # 0
arr.append(0)
print(arr)
arr[0] = (1 << 32) - 1
print(arr[0])  # 2147483648

print(1 << 3 - 1)
