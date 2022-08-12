
            self._tree[rt] = self._tree[rt * 2]
        if self._tree[rt * 2 + 1] > self._tree[rt]:
            self._tree[rt] = self._tree[rt * 2 + 1