如果是单链表，我们无法在 $O(1)$ 的时间拿到前驱节点，这也是为什么我们遍历的时候老是维护一个前驱节点的原因。

如果链表的头节点可能会被删除或者发生变化 我们可以引入 dummy 头节点
链表的基本操作

开放的题可以类比数组

**插入**
temp = 待插入位置的前驱节点.next
待插入位置的前驱节点.next = 待插入指针
待插入指针.next = temp
**删除**
pre.next=pre.next.next

```JS
 while (dummyP && dummyP.next) {
    let next: Node | undefined = dummyP.next
    if (next.value === val) {
      dummyP.next = next.next
      next = next.next
    }
    // 只有下个节点不是要删除的节点才更新 current
    if (!next || next.value !== val) dummyP = next
  }
```

**遍历**
当前指针 = 头指针
while 当前节点不为空 {
print(当前节点)
当前指针 = 当前指针.next
}

**翻转整段链表**
tmp=p2.next
p2.next=p1
p1=p2
p2=tmp

k 个一组翻转链表
