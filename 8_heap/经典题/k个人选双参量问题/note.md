<!-- 双变量制约 -->

2 个变量，想办法”定“住一个:提前把一个排好序(涉及到最值的变量) 取的时候已经是最小了
总结:`乘参量 1 用排序`，`加参量 2 用堆维护`,出堆入堆形成抗衡，同时更新 res

`857. 雇佣 K 名工人的最低成本.py`
题目中对 wage/quality 升序排列，对 quality 用堆降序维护。每次入堆 wage/quality 不断变大，出堆 quality 不断变小,出堆入堆形成抗衡

`1383. 最大的团队表现值.py`
题目中 对 eff 降序排列，堆 speed 升序维护:每次入堆 eff 最没用的，出堆 speed 最没用的,出堆入堆形成抗衡

即：`进来的越没用，出去的最没用`

---

两条原则:

- 2 个变量，想办法”定“住一个 => `排序`
- 进来一个 x 维度上最好的，出去一个 y 维度上最没用的
