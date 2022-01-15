class Node {
  constructor(val) {
    this.val = val
    this.cnt = 1
    this.size = 1
    this.fac = Math.random()
    this.left = null
    this.right = null
  }
  push_up() {
    let tmp = this.cnt
    tmp += Node.getSize(this.left)
    tmp += Node.getSize(this.right)
    this.size = tmp
  }
  rotate_right() {
    let node = this
    let left = node.left
    node.left = left.right
    left.right = node
    node = left
    node.right.push_up()
    node.push_up()
    return node
  }
  rotate_left() {
    let node = this
    let right = node.right
    node.right = right.left
    right.left = node
    node = right
    node.left.push_up()
    node.push_up()
    return node
  }

  static getSize(node) {
    // return node?.size || 0;
    return (node && node.size) || 0
  }

  static getFac(node) {
    return (node && node.fac) || 0
  }
}
class Treap {
  constructor(compare, left = -Infinity, right = Infinity) {
    this.root = new Node(right)
    this.root.fac = Infinity
    this.root.left = new Node(left)
    this.root.left.fac = -Infinity
    this.root.push_up()

    this.compare = compare
  }
  get size() {
    return this.root.size - 2
  }
  get height() {
    function getHeight(node) {
      if (node === null) return 0
      return 1 + Math.max(getHeight(node.left), getHeight(node.right))
    }
    return getHeight(this.root)
  }
  insert(val) {
    let compare = this.compare
    // js 里没 & 这种引用  所以要带着parent和上次的方向  在c++里直接 Tree &rt 就可以了
    function dfs(node, val, parent, direction) {
      if (compare(node.val, val) === 0) {
        node.cnt++
        node.push_up()
      } else if (compare(node.val, val) === 1) {
        if (node.left) {
          dfs(node.left, val, node, 'left')
        } else {
          node.left = new Node(val)
          node.push_up()
        }

        if (Node.getFac(node.left) > node.fac) {
          parent[direction] = node.rotate_right()
        }
      } else if (compare(node.val, val) === -1) {
        if (node.right) {
          dfs(node.right, val, node, 'right')
        } else {
          node.right = new Node(val)
          node.push_up()
        }

        if (Node.getFac(node.right) > node.fac) {
          parent[direction] = node.rotate_left()
        }
      }
      parent.push_up()
    }
    dfs(this.root.left, val, this.root, 'left')
  }
  remove(val) {
    let compare = this.compare
    function dfs(node, val, parent, direction) {
      if (node === null) return

      if (compare(node.val, val) === 0) {
        if (node.cnt > 1) {
          node.cnt--
          node.push_up()
        } else if (node.left === null && node.right === null) {
          parent[direction] = null
        } else {
          // 旋到根节点
          if (node.right === null || Node.getFac(node.left) > Node.getFac(node.right)) {
            parent[direction] = node.rotate_right()
            dfs(parent[direction].right, val, parent[direction], 'right')
          } else {
            parent[direction] = node.rotate_left()
            dfs(parent[direction].left, val, parent[direction], 'left')
          }
        }
      } else if (compare(node.val, val) === 1) {
        dfs(node.left, val, node, 'left')
      } else if (compare(node.val, val) === -1) {
        dfs(node.right, val, node, 'right')
      }
      parent.push_up()
    }
    dfs(this.root.left, val, this.root, 'left')
  }

  getRankByVal(val) {
    if (val === void 0) return 0
    let compare = this.compare
    function dfs(node, val) {
      if (node === null) return 0

      if (compare(node.val, val) === 0) {
        // return ((node.left && node.left.size) || 0) + 1;
        return Node.getSize(node.left) + 1
      } else if (compare(node.val, val) === 1) {
        return dfs(node.left, val)
      } else if (compare(node.val, val) === -1) {
        // return dfs(node.right, val) + ((node.left && node.left.size) || 0) + node.cnt;
        return dfs(node.right, val) + Node.getSize(node.left) + node.cnt
      }
    }
    // 因为有个-Infinity 所以-1
    return dfs(this.root, val) - 1
  }
  getValByRank(rank) {
    if (rank === void 0) return Infinity
    function dfs(node, rank) {
      if (node === null) return Infinity

      if (Node.getSize(node.left) >= rank) {
        return dfs(node.left, rank)
      } else if (Node.getSize(node.left) + node.cnt >= rank) {
        return node.val
      } else {
        return dfs(node.right, rank - Node.getSize(node.left) - node.cnt)
      }
    }
    // 因为有个-Infinity 所以 + 1
    return dfs(this.root, rank + 1)
  }
  // lower_bound - 1
  getPrev(val) {
    if (val === void 0) return -Infinity
    let compare = this.compare
    function dfs(node, val) {
      if (node === null) return -Infinity
      if (compare(node.val, val) >= 0) return dfs(node.left, val)

      let tmp = dfs(node.right, val)
      if (compare(node.val, tmp) == 1) {
        return node.val
      } else {
        return tmp
      }
    }
    return dfs(this.root, val)
  }
  // upper_bound
  getNext(val) {
    if (val === void 0) return Infinity
    let compare = this.compare
    function dfs(node, val) {
      if (node === null) return Infinity
      if (compare(node.val, val) <= 0) return dfs(node.right, val)

      let tmp = dfs(node.left, val)
      if (compare(node.val, tmp) < 0) {
        return node.val
      } else {
        return tmp
      }
    }
    return dfs(this.root, val)
  }
  // 小于等于
  getPrevPlus(val) {
    if (val === void 0) return -Infinity
    let compare = this.compare
    function dfs(node, val) {
      if (node === null) return -Infinity
      if (compare(node.val, val) === 0) return node.val
      if (compare(node.val, val) >= 0) return dfs(node.left, val)

      let tmp = dfs(node.right, val)
      if (compare(node.val, tmp) == 1) {
        return node.val
      } else {
        return tmp
      }
    }
    return dfs(this.root, val)
  }
  // 大于等于
  getNextPlus(val) {
    if (val === void 0) return Infinity
    let compare = this.compare
    function dfs(node, val) {
      if (node === null) return Infinity
      if (compare(node.val, val) === 0) return node.val
      if (compare(node.val, val) < 0) return dfs(node.right, val)

      let tmp = dfs(node.left, val)
      if (compare(node.val, tmp) < 0) {
        return node.val
      } else {
        return tmp
      }
    }
    return dfs(this.root, val)
  }
  find(val) {
    if (val === void 0) return null
    let compare = this.compare
    function dfs(node, val) {
      if (node === null) return null
      if (compare(node.val, val) === 0) return node.val
      if (compare(node.val, val) < 0) return dfs(node.right, val)
      return dfs(node.left, val)
    }
    return dfs(this.root, val)
  }
}

function isNStraightHand(hand, groupSize) {
  let set = new Treap((node, val) => {
    if (node == val) return 0
    if (node > val) return 1
    return -1
  })

  for (let i = 0; i < hand.length; i++) {
    set.insert(hand[i])
  }

  while (set.size) {
    let head = set.getNext(-Infinity)
    if (head == Infinity) return false
    for (let i = 0; i < groupSize; i++) {
      if (set.find(head + i) == null) {
        return false
      } else {
        set.remove(head + i)
      }
    }
  }
  return true
}

if (require.main === module) {
  const set = new Treap((node, val) => {
    if (node == val) return 0
    if (node > val) return 1
    return -1
  })
  set.insert(1)
  set.insert(2)
  set.insert(3)
  set.insert(4)
  console.dir(set, { depth: null })
  console.log(set.getNextPlus(2.1))
}
