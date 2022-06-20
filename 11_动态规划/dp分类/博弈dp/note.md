记忆化搜索容易理解

多以 dfs(index)为搜索函数
`能使对方输就是自己赢，不能使对方输就是自己输`
`自己赢=自己赢+对手不赢`

```Python
mem = {}
def win(A):                     # 判断状态A是否为胜态
    if A not in mem:
        if is_final(A):         # 若A为终局态
            mem[A] = rule(A)    # 根据游戏规则判断A的胜负
        else:                   # 若A为非终局态，则根据策梅洛定理判断其胜负
            mem[A] = not all(win(B) for B in next_states(A))
                                # next_states(A)返回A的所有次态
    return mem[A]
```
