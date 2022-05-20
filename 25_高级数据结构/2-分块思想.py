# // sqrt_decomposition
# 完整段：查询用sum，修改改sum、标记这个块加了多少add
# 不完整快：查询用nums[i]+add ，修改改nums[i]和sum
from math import sqrt
import sys


def getChunkId(index: int) -> int:
    return index // chunkLen


def update(left: int, right: int, delta: int) -> None:
    """
    完整段：chunkAdd chunkSum 
    不完整段：nums[i] chunkSum
    """
    lChunk, rChunk = getChunkId(left), getChunkId(right)
    if lChunk == rChunk:  # 段内直接暴力
        for i in range(left, right + 1):
            nums[i] += delta
            chunkSum[lChunk] += delta
    else:
        i, j = left, right
        # 暴力修改两边
        while getChunkId(i) == lChunk:
            nums[i] += delta
            chunkSum[lChunk] += delta
            i += 1
        while getChunkId(j) == rChunk:
            nums[j] += delta
            chunkSum[rChunk] += delta
            j -= 1
        # 完整段 用chunkAdd代替每个点的实际修改
        for i in range(lChunk + 1, rChunk):
            chunkAdd[i] += delta
            chunkSum[i] += delta * chunkLen


def query(left: int, right: int) -> int:
    """
    完整段：chunkSum 
    不完整段：nums[i]+chunkAdd
    """
    lChunk, rChunk = getChunkId(left), getChunkId(right)
    if lChunk == rChunk:
        return sum(nums[left : right + 1]) + chunkAdd[lChunk] * (right - left + 1)
    else:
        i, j = left, right
        res = 0
        # 暴力修改两边
        while getChunkId(i) == lChunk:
            res += nums[i] + chunkAdd[lChunk]
            i += 1
        while getChunkId(j) == rChunk:
            res += nums[j] + chunkAdd[rChunk]
            j -= 1
        # 完整段
        for i in range(lChunk + 1, rChunk):
            res += chunkSum[i]
        return res


input = sys.stdin.readline

n, m = map(int, input().split())
nums = list(map(int, input().split()))

chunkLen = int(sqrt(n))  # 每段长度
chunkAdd = [0] * (n // chunkLen + 1)
chunkSum = [0] * (n // chunkLen + 1)

for i in range(n):
    chunkSum[getChunkId(i)] += nums[i]

for _ in range(m):
    opt = input().split()
    if opt[0] == 'C':
        left, right, delta = map(int, opt[1:])
        left, right = left - 1, right - 1
        update(left, right, delta)
    else:
        left, right = map(int, opt[1:])
        left, right = left - 1, right - 1
        print(query(left, right))
