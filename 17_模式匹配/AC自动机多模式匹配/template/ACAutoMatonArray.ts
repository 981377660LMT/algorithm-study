// TODO: to数组压缩成一维
class ACAutoMatonArray {
  words: number[] // words[i] 表示加入的第i个模式串对应的节点编号.
  parent: number[] // parent[v] 表示节点v的父节点.
  sigma: number // 字符集大小.
  offset: number // 字符集的偏移量.
  children: number[][] // children[v][c] 表示节点v通过字符c转移到的节点.
  suffixLink: number[] // 又叫fail.指向当前节点最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
  bfsOrder: number[] // 结点的拓扑序,0表示虚拟节点.
  needUpdateChildren: boolean // 是否需要更新children数组.

  constructor(sigma: number, offset: number) {
    this.sigma = sigma
    this.offset = offset
    this.newNode()
  }

  addString(str: string): number {
    if (str.length === 0) {
      return 0
    }
    let pos = 0
    for (let s of str) {
      let ord = s.charCodeAt(0) - this.offset
      if (this.children[pos][ord] === -1) {
        this.children[pos][ord] = this.newNode()
        this.parent.push(pos)
      }
      pos = this.children[pos][ord]
    }
    this.words.push(pos)
    return pos
  }

  addChar(pos: number, ord: number): number {
    ord -= this.offset
    if (this.children[pos][ord] !== -1) {
      return this.children[pos][ord]
    }
    this.children[pos][ord] = this.newNode()
    this.parent.push(pos)
    return this.children[pos][ord]
  }

  move(pos: number, ord: number): number {
    ord -= this.offset
    if (this.needUpdateChildren) {
      return this.children[pos][ord]
    }
    while (true) {
      let nexts = this.children[pos]
      if (nexts[ord] !== -1) {
        return nexts[ord]
      }
      if (pos === 0) {
        return 0
      }
      pos = this.suffixLink[pos]
    }
  }

  size(): number {
    return this.children.length
  }

  buildSuffixLink(needUpdateChildren: boolean) {
    this.needUpdateChildren = needUpdateChildren
    this.suffixLink = new Array(this.children.length).fill(-1)
    this.bfsOrder = new Array(this.children.length).fill(0)
    let head = 0,
      tail = 0
    this.bfsOrder[tail] = 0
    tail++
    while (head < tail) {
      let v = this.bfsOrder[head]
      head++
      for (let i = 0; i < this.children[v].length; i++) {
        let next = this.children[v][i]
        if (next === -1) {
          continue
        }
        this.bfsOrder[tail] = next
        tail++
        let f = this.suffixLink[v]
        while (f !== -1 && this.children[f][i] === -1) {
          f = this.suffixLink[f]
        }
        this.suffixLink[next] = f
        if (f === -1) {
          this.suffixLink[next] = 0
        } else {
          this.suffixLink[next] = this.children[f][i]
        }
      }
    }
    if (!needUpdateChildren) {
      return
    }
    for (let v of this.bfsOrder) {
      for (let i = 0; i < this.children[v].length; i++) {
        let next = this.children[v][i]
        if (next === -1) {
          let f = this.suffixLink[v]
          if (f === -1) {
            this.children[v][i] = 0
          } else {
            this.children[v][i] = this.children[f][i]
          }
        }
      }
    }
  }

  getCounter(): number[] {
    let counter = new Array(this.children.length).fill(0)
    for (let pos of this.words) {
      counter[pos]++
    }
    for (let v of this.bfsOrder) {
      if (v !== 0) {
        counter[v] += counter[this.suffixLink[v]]
      }
    }
    return counter
  }

  getIndexes(): number[][] {
    let res = new Array(this.children.length).fill([])
    for (let i = 0; i < this.words.length; i++) {
      let pos = this.words[i]
      res[pos].push(i)
    }
    for (let v of this.bfsOrder) {
      if (v !== 0) {
        let from = this.suffixLink[v],
          to = v
        let arr1 = res[from],
          arr2 = res[to]
        let arr3 = []
        let i = 0,
          j = 0
        while (i < arr1.length && j < arr2.length) {
          while (i < arr1.length && j < arr2.length) {
            if (arr1[i] < arr2[j]) {
              arr3.push(arr1[i])
              i++
            } else if (arr1[i] > arr2[j]) {
              arr3.push(arr2[j])
              j++
            } else {
              arr3.push(arr1[i])
              i++
              j++
            }
          }
        }
        while (i < arr1.length) {
          arr3.push(arr1[i])
          i++
        }
        while (j < arr2.length) {
          arr3.push(arr2[j])
          j++
        }
        res[to] = arr3
      }
    }
    return res
  }

  dp(f: (from: number, to: number) => void) {
    for (let v of this.bfsOrder) {
      if (v !== 0) {
        f(this.suffixLink[v], v)
      }
    }
  }

  newNode(): number {
    this.parent.push(-1)
    let nexts = new Array(this.sigma).fill(-1)
    this.children.push(nexts)
    return this.children.length - 1
  }
}
