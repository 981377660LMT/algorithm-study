// 跳表是在 O(log(n)) 时间内完成增加、删除、搜索操作的数据结构。
// 跳表相比于树堆与红黑树，其功能与性能相当，并且跳表的代码长度相较下更短，
// 其设计思想与链表相似。
// 在Redis中的有序集合（Sorted Set）就是用跳表来实现的
// 跳表支持区间查询
import assert from 'assert'

class SkipNode {
  public value: number
  public right: SkipNode | null
  public down: SkipNode | null

  constructor(value: number) {
    this.value = value
    this.right = null
    this.down = null
  }
}

interface ISkipList {
  search: (target: number) => boolean
  add: (num: number) => void
  erase: (num: number) => boolean
}

class SkipList implements ISkipList {
  private head: SkipNode
  private maxLevel: number

  constructor(maxLevel: number = 16) {
    const leftWall = Array.from({ length: maxLevel }, () => new SkipNode(-Infinity))
    const rightWall = Array.from({ length: maxLevel }, () => new SkipNode(Infinity))
    for (let i = 0; i < maxLevel - 1; i++) {
      leftWall[i].right = rightWall[i]
      leftWall[i].down = leftWall[i + 1]
      rightWall[i].down = rightWall[i + 1]
    }
    leftWall[maxLevel - 1].right = rightWall[maxLevel - 1]
    this.head = leftWall[0]
    this.maxLevel = maxLevel
  }

  search(target: number): boolean {
    let headP: SkipNode | null = this.head
    while (headP) {
      if (headP.right!.value > target) {
        headP = headP.down
      } else if (headP.right!.value < target) {
        headP = headP.right
      } else {
        return true
      }
    }

    return false
  }

  /**
   * @description num一个一个看，插入
   * @description 随机层数，如果遇到节点右边的值比num大，则记录下来，等待之后插入到其左边;否则向右移动
   * @description 左右插入普通方法即可，上下穿针需要dummy节点
   */
  add(num: number): void {
    let level = 0
    // 每一层的pre，即待插入在这个节点后面
    const preNodes = Array<SkipNode>(this.maxLevel)
    let headP: SkipNode | null = this.head
    while (headP) {
      if (headP.right!.value >= num) {
        preNodes[level] = headP
        headP = headP.down
        level++
      } else {
        headP = headP.right
      }
    }
    console.log(preNodes)
    let dummy = new SkipNode(0)
    const nodesToInsert = Array.from({ length: this.getRandomLevel() }, () => new SkipNode(num))
    const insertNum = nodesToInsert.length
    // 注意层数是this.maxLevel - insertNum 一般很大在下面
    for (let i = 0, j = this.maxLevel - insertNum; i < insertNum; i++, j++) {
      const insertNode = nodesToInsert[i]
      const preNode = preNodes[j]
      // 水平方向
      insertNode.right = preNode.right
      preNode.right = insertNode
      // 竖直方向
      dummy.down = insertNode
      dummy = insertNode
    }
  }

  erase(num: number): boolean {
    let res = false
    let headP: SkipNode | null = this.head
    while (headP) {
      if (headP.right!.value > num) {
        headP = headP.down
      } else if (headP.right!.value < num) {
        headP = headP.right
      } else {
        res = true
        headP.right = headP.right?.right || null
        // 这里是为什么
        headP = headP.down
      }
    }

    return res
  }

  private getRandomLevel(): number {
    const maxRand = 2 ** this.maxLevel - 1
    // 跳表需要对数概率 randomLevel一般很小
    const randomLevel = this.maxLevel - Math.floor(Math.log2(1 + Math.random() * maxRand))
    return randomLevel
  }
}

if (require.main === module) {
  const skiplist = new SkipList()
  skiplist.add(1)
  skiplist.add(2)
  skiplist.add(3)
  assert.strictEqual(skiplist.search(0), false, '返回 false，0 不在跳表中')
  skiplist.add(4)
  assert.strictEqual(skiplist.search(1), true)
  assert.strictEqual(skiplist.erase(0), false, '返回 false，0 不在跳表中')
  assert.strictEqual(skiplist.erase(1), true)
  assert.strictEqual(skiplist.search(1), false, '返回 false，1 已被擦除')
}

export { SkipList }
