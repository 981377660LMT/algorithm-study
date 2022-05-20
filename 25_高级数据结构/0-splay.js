// https://www.acwing.com/solution/content/81783/
class Node {
  constructor(v, p) {
    this.v = v
    this.s = [null, null]
    this.p = p

    this.size = 1
    this.flag = 0
  }

  static getSize(node) {
    return (node && node.size) || 0
  }

  static getFac(node) {
    return (node && node.fac) || 0
  }
}

class Splay {
  constructor(arr) {
    this.arr = arr
    this.root = 0
    this.tr = {}
    this.idx = 0

    for (let i = 0; i <= arr.length + 1; i++) {
      this.insert(i)
    }
  }

  push_up(x) {
    const tr = this.tr
    // tr[x].size = tr[tr[x].s[0]].size + tr[tr[x].s[1]].size + 1;
    tr[x].size = Node.getSize(tr[tr[x].s[0]]) + Node.getSize(tr[tr[x].s[1]]) + 1
  }

  push_down(x) {
    const tr = this.tr
    if (tr[x].flag) {
      // 调换位置并将标记下传
      ;[tr[x].s[0], tr[x].s[1]] = [tr[x].s[1], tr[x].s[0]]
      if (tr[x].s[0]) {
        tr[tr[x].s[0]].flag ^= 1
      }
      if (tr[x].s[1]) {
        tr[tr[x].s[1]].flag ^= 1
      }
      tr[x].flag = 0
    }
  }

  // 转谁 谁上去
  rotate(x) {
    const tr = this.tr
    let y = tr[x].p
    let z = tr[y].p

    let k = +(tr[y].s[1] == x)
    // 1 右子 0 左子

    // 上去
    if (tr[z]) {
      tr[z].s[+(tr[z].s[1] == y)] = x
    }
    tr[x].p = z

    // 被顶下来的y  用原来指向x的指针接好自己新位置的原住民
    tr[y].s[k] = tr[x].s[k ^ 1]
    if (tr[tr[x].s[k ^ 1]]) {
      tr[tr[x].s[k ^ 1]].p = y
    }

    // x y 建立联系
    tr[x].s[k ^ 1] = y
    tr[y].p = x

    // 维护每个点的size
    this.push_up(y)
    this.push_up(x)
  }

  splay(x, k) {
    const tr = this.tr
    while (tr[x].p != k) {
      let y = tr[x].p
      let z = tr[y].p
      if (z !== k) {
        if (+(tr[y].s[1] == x) ^ +(tr[z].s[1] == y)) {
          // 如果是折线
          this.rotate(x)
        } else {
          this.rotate(y)
        }
      }
      this.rotate(x)
    }
    // k == 0 转到根下
    if (!k) {
      this.root = x
    }
  }

  insert(v) {
    const tr = this.tr
    let u = this.root
    let p = 0
    while (u) {
      p = u
      u = tr[u].s[+(v > tr[u].v)]
    }
    // 这样就是 [l,r] l  r + 2了
    u = ++this.idx
    if (tr[p]) {
      tr[p].s[+(v > tr[p].v)] = u
    }
    tr[u] = new Node(v, p)
    this.splay(u, 0)
  }

  get_k(k) {
    const tr = this.tr
    let u = this.root
    while (true) {
      this.push_down(u)
      if (Node.getSize(tr[tr[u].s[0]]) >= k) {
        u = tr[u].s[0]
      } else if (Node.getSize(tr[tr[u].s[0]]) + 1 == k) {
        return u
      } else {
        k -= Node.getSize(tr[tr[u].s[0]]) + 1
        u = tr[u].s[1]
      }
    }
    return -1
  }

  output(u) {
    const tr = this.tr
    const push_down = this.push_down.bind(this)
    const arr = this.arr

    let ret = []
    function dfs(u) {
      push_down(u)

      if (tr[u].s[0]) {
        dfs(tr[u].s[0])
      }
      // 不是哨兵
      if (tr[u].v >= 1 && tr[u].v <= arr.length) {
        ret.push(tr[u].v)
      }
      if (tr[u].s[1]) {
        dfs(tr[u].s[1])
      }
    }
    dfs(this.root)
    return ret
  }
}

var fs = require('fs')
var buf = ''
process.stdin.on('readable', function () {
  var chunk = process.stdin.read()
  if (chunk) buf += chunk.toString()
})
process.stdin.on('end', function () {
  let ret = 0
  let m
  let n
  let splay
  buf.split('\n').forEach(function (line, index) {
    let num = +line
    if (index === 0) {
      tmp = line.split(' ').filter(item => !!item)
      n = +tmp[0]
      m = +tmp[1]
      splay = new Splay(new Array(n))
    } else if (index <= m) {
      tmp = line.split(' ').filter(item => !!item)
      let l = splay.get_k(+tmp[0])
      let r = splay.get_k(+tmp[1] + 2)
      splay.splay(l, 0)
      splay.splay(r, l)
      splay.tr[splay.tr[r].s[0]].flag ^= 1
    }
  })

  console.log(splay.output().join(' '))
})
