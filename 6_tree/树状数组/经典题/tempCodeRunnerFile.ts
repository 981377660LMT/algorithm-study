{
      tree[i] = max[i]
      for (let j = 1; j < (i & -i); j <<= 1) tree[i] = Math.max(tree[i], tree[i - j])
    }