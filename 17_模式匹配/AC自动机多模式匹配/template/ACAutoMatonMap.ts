/* eslint-disable no-loop-func */
/* eslint-disable max-len */

class ACAutoMatonMap {
  /** wordPos[i] 表示加入的第i个模式串对应的节点编号. */
  readonly wordPos: number[] = [] //
  private readonly _children: Map<number, number>[] = [new Map()]
  private _suffixLink!: Int32Array
  private _bfsOrder!: Int32Array

  addString(str: string): number {
    if (str.length === 0) return 0
    let pos = 0
    for (let i = 0; i < str.length; i++) {
      const ord = str[i].charCodeAt(0)
      const nexts = this._children[pos]
      if (nexts.has(ord)) {
        pos = nexts.get(ord)!
      } else {
        const nextState = this._children.length
        nexts.set(ord, nextState)
        pos = nextState
        this._children.push(new Map())
      }
    }
    this.wordPos.push(pos)
    return pos
  }

  addChar(pos: number, ord: number): number {
    let nexts = this._children[pos]
    if (nexts.has(ord)) {
      return nexts.get(ord)!
    }
    const nextState = this._children.length
    nexts.set(ord, nextState)
    this._children.push(new Map())
    return nextState
  }

  move(pos: number, ord: number): number {
    while (true) {
      const nexts = this._children[pos]
      if (nexts.has(ord)) {
        return nexts.get(ord)!
      }
      if (pos === 0) {
        return 0
      }
      pos = this._suffixLink[pos]
    }
  }

  buildSuffixLink() {
    this._suffixLink = new Int32Array(this._children.length).fill(-1)
    this._bfsOrder = new Int32Array(this._children.length)
    let head = 0
    let tail = 1
    while (head < tail) {
      const v = this._bfsOrder[head]
      head++
      this._children[v].forEach((next, char) => {
        this._bfsOrder[tail] = next
        tail++
        let f = this._suffixLink[v]
        while (f !== -1) {
          if (this._children[f].has(char)) {
            break
          }
          f = this._suffixLink[f]
        }
        if (f === -1) {
          this._suffixLink[next] = 0
        } else {
          this._suffixLink[next] = this._children[f].get(char)!
        }
      })
    }
  }

  getCounter(): Uint32Array {
    const counter = new Uint32Array(this._children.length)
    for (let i=
    for (let v of this._bfsOrder) {
      if (v !== 0) {
        counter[v] += counter[this._suffixLink[v]]
      }
    }
    return counter
  }

  getIndexes(): number[][] {
    let res = new Array(this._children.length).fill([])
    for (let i = 0; i < this.wordPos.length; i++) {
      let pos = this.wordPos[i]
      res[pos].push(i)
    }
    for (let v of this._bfsOrder) {
      if (v !== 0) {
        let from = this._suffixLink[v],
          to = v
        let arr1 = res[from],
          arr2 = res[to]
        let arr3 = []
        let i = 0,
          j = 0
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
    for (let v of this._bfsOrder) {
      if (v !== 0) {
        f(this._suffixLink[v], v)
      }
    }
  }

  size(): number {
    return this._children.length
  }
}
