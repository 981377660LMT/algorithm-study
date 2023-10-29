class Solution:
    def sumCounts(self, nums: List[int]) -> int:
        Mod = int(1e9 + 7)

        def push_lazy(node, left, right):
            tmp = lazy[node]
            if lazy[node] != 0:
                he[node] += tmp * tmp * (right - left + 1) + tree[node] * 2 * tmp
                tree[node] += tmp * (right - left + 1)
                if left != right:
                    lazy[node * 2] += tmp
                    lazy[node * 2 + 1] += tmp
                lazy[node] = 0

        def update(node, left, right, l, r, val):
            push_lazy(node, left, right)
            if r < left or l > right:
                return
            if l <= left and r >= right:
                lazy[node] += val
                push_lazy(node, left, right)
                return
            mid = (left + right) // 2
            update(node * 2, left, mid, l, r, val)
            update(node * 2 + 1, mid + 1, right, l, r, val)
            tree[node] = tree[node * 2] + tree[node * 2 + 1]
            he[node] = he[node * 2] + he[node * 2 + 1]

        def query2(node, left, right, l, r):
            push_lazy(node, left, right)
            if r < left or l > right:
                return 0
            if l <= left and r >= right:
                return he[node]
            mid = (left + right) // 2
            sum_left = query2(node * 2, left, mid, l, r)
            sum_right = query2(node * 2 + 1, mid + 1, right, l, r)
            return sum_left + sum_right

        n = len(nums)
        tree = [0] * (4 * n)
        lazy = [0] * (4 * n)
        he = [0] * (4 * n)

        s = set()
        d = defaultdict(lambda: n)
        right = [n] * n
        for i in range(n - 1, -1, -1):
            right[i] = d[nums[i]]
            d[nums[i]] = i

        for i, v in enumerate(nums):
            s.add(v)
            update(1, 0, n - 1, i, i, len(s))

        ans = 0
        for i in range(n):
            ans += query2(1, 0, n - 1, i, n - 1)
            update(1, 0, n - 1, i, right[i] - 1, -1)
        return ans % Mod
