## 模拟一个过程，而这个过程一般是按照 `时间顺序` 去执行:

0. 分析状态与每种状态的优先级，用 SortedList/heap 存储；维护一个时间戳 curTime。
1. 弄清模拟的结束条件 (while ...)
   ```py
   while pq or ei < len(events): # !所有的会议处理完毕，结束循环
   while len(res) < n:
   while remain > 0 or rightFinish or rightWait:  # !当旧仓库还有货物或者右边还有人要回来时
   ```
2. 每次 while 循环处理事件 :
   - while 一次性加入所有开始时间小于等于当前时间的 event/删除所有结束时间小于当前时间的 event
   - if 来处理不同 event
     **注意如果没有 event 要处理，则需要更新时间到下一次状态变化的时间**
