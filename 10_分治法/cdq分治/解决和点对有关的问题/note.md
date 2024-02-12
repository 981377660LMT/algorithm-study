/**
if problems have following characteristics-
1.it can be solved offline
2.if every insert operations were before query operations
each query operations could have been solved in O(L)(i mean faster like logn)
3.query operation is cumulative i.e. we don't need every insert at once to solve this query
this type of problem can be solved using divide and conquer
complexity: nlogn\*O(L)
**/
Problem:

    You have an empty set. You need to perform following operations –
    Insert a given number X in the set.
    Count how many numbers in the set is less than or equal to a given number X

Solution:

    op[1...n] = operations
    ans[1...n] = array to store answers
    solve(l, r) {
      if op[l...r] has no queries: return
      m = (l + r) / 2
      solve(l, m)
      solve(m+1, r)
      ds = statically built DS using insert operations in op[l...m]
      for each query operation i in op[m+1...r]:
        ans[i] += ds.query(op[i])
    }

# 偏序问题是一个挺老生常谈的问题。

基本上就是说某一种东西有那么 2~3 种属性，然后给你两个这个东西，当他们的每种属性之间都必须满足一些给定的大小关系时，约束条件成立。

举两个例子，一个是逆序对问题，一个是最长上升子序列问题。这些都是最典型的“偏序约束条件”。
