import sys
import threading


def main():
    import sys

    sys.setrecursionlimit(1 << 25)

    N, Q = map(int, sys.stdin.readline().split())
    H = list(map(int, sys.stdin.readline().split()))
    queries = []
    for idx in range(Q):
        l_i, r_i = map(int, sys.stdin.readline().split())
        queries.append((r_i, l_i, idx))

    # Compute left_taller for each building
    left_taller = [0] * (N + 2)  # 1-based indexing
    stack = []
    for k in range(1, N + 1):
        while stack and H[stack[-1] - 1] < H[k - 1]:
            stack.pop()
        if not stack:
            left_taller[k] = 0
        else:
            left_taller[k] = stack[-1]
        stack.append(k)

    # Prepare queries
    # queries: list of (ri, li, idx)
    # We need to process queries with larger ri first
    queries.sort(reverse=True)

    # Initialize BIT
    class BIT:
        def __init__(self, size):
            self.N = size
            self.tree = [0] * (self.N + 2)

        def update(self, i, delta):
            while i <= self.N:
                self.tree[i] += delta
                i += i & -i

        def query(self, i):
            res = 0
            while i > 0:
                res += self.tree[i]
                i -= i & -i
            return res

    bit_size = N + 2  # Since positions can be up to N + 1
    bit = BIT(bit_size)

    answers = [0] * Q
    query_idx = 0
    k = N + 1
    while k >= 1:
        k -= 1
        # Insert left_taller[k] into BIT
        if k >= 1:
            position = left_taller[k] + 1  # Convert to 1-based indexing
            bit.update(position, 1)
        # Process queries with ri == k
        while query_idx < Q and queries[query_idx][0] == k:
            ri, li, idx = queries[query_idx]
            # Query BIT for positions less than li
            count = bit.query(li - 1)
            answers[idx] = count
            query_idx += 1

    # Since we sorted the queries, we need to rearrange the answers
    # But we stored the index with each query, so we can output directly
    for ans in answers:
        print(ans)


threading.Thread(target=main).start()
