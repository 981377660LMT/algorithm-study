from collections import defaultdict


trie = lambda: defaultdict(trie)  # type: ignore
trie = trie()
