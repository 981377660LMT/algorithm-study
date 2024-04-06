https://notes.sshwy.name/Randomized/Count-min-Sketch/
「随机算法专题」Count-min Sketch 入门
这是一种，估计出现次数的数据结构，适用于数据流（强在线），不同的元素数不多的情况。

基本思想是一个哈希表可以估计某个值的出现的次数，但是会估多，因此多个哈希取最小值即可。这里的哈希要取均匀哈希才能保证理论复杂度的正确性。
