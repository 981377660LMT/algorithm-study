 TrieNode()
    t.children["a"] = TrieNode()
    t.children["a"].children["b"] = TrieNode()
    print(Solution().serialize(t))
