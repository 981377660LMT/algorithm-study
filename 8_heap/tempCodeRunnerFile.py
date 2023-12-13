import sys

input = sys.stdin.buffer.readline
n, q = map(int, input().split())
s = list(map(int, input().split()))
heap = MinMaxHeap(s)
for _ in range(q):
    query = tuple(map(int, input().split()))
    if query[0] == 0:
        heap.push(query[1])
    elif query[0] == 1:
        print(heap.popMin())
    else:
        print(heap.popMax())
