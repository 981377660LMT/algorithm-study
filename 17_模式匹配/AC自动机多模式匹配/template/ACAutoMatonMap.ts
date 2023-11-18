const INF = 2e18

class ACAutoMatonMap {
  words: number[]
  bfsOrder: number[]
  children: Map<string, number>[]
  suffixLink: number[]

  constructor() {
    this.words = []
    this.bfsOrder = []
    this.children = [new Map()]
    this.suffixLink = []
  }

  addString(str: string): number {
    if (!str) {
      return 0
    }
    let pos = 0
    for (const char of str) {
      const nexts = this.children[pos]
      if (nexts.has(char)) {
        pos = nexts.get(char)!
      } else {
        const nextState = this.children.length
        nexts.set(char, nextState)
        pos = nextState
        this.children.push(new Map())
      }
    }
    this.words.push(pos)
    return pos
  }

  addChar(pos: number, char: string): number {
    const nexts = this.children[pos]
    if (nexts.has(char)) {
      return nexts.get(char)!
    }
    const nextState = this.children.length
    nexts.set(char, nextState)
    this.children.push(new Map())
    return nextState
  }

  move(pos: number, char: string): number {
    while (true) {
      const nexts = this.children[pos]
      if (nexts.has(char)) {
        return nexts.get(char)!
      }
      if (pos === 0) {
        return 0
      }
      pos = this.suffixLink[pos]
    }
  }

  buildSuffixLink(): void {
    this.suffixLink = new Array(this.children.length).fill(-1)
    this.bfsOrder = new Array(this.children.length).fill(0)
    let head = 0,
      tail = 1
    while (head < tail) {
      const v = this.bfsOrder[head++]
      for (const [char, next] of this.children[v]) {
        this.bfsOrder[tail++] = next
        let f = this.suffixLink[v]
        while (f !== -1 && !this.children[f].has(char)) {
          f = this.suffixLink[f]
        }
        this.suffixLink[next] = f === -1 ? 0 : this.children[f].get(char)!
      }
    }
  }

  getCounter(): number[] {
    const counter = new Array(this.children.length).fill(0)
    for (const pos of this.words) {
      counter[pos]++
    }
    for (const v of this.bfsOrder) {
      if (v !== 0) {
        counter[v] += counter[this.suffixLink[v]]
      }
    }
    return counter
  }

  getIndexes(): number[][] {
    const res = Array.from({ length: this.children.length }, () => [])
    for (const [i, pos] of this.words.entries()) {
      res[pos].push(i)
    }
    for (const v of this.bfsOrder) {
      if (v !== 0) {
        const from = this.suffixLink[v]
        const arr1 = res[from],
          arr2 = res[v]
        const arr3 = []
        let i = 0,
          j = 0
        while (i < arr1.length && j < arr2.length) {
          if (arr1[i] < arr2[j]) {
            arr3.push(arr1[i++])
          } else if (arr1[i] > arr2[j]) {
            arr3.push(arr2[j++])
          } else {
            arr3.push(arr1[i++])
            j++
          }
        }
        while (i < arr1.length) {
          arr3.push(arr1[i++])
        }
        while (j < arr2.length) {
          arr3.push(arr2[j++])
        }
        res[v] = arr3
      }
    }
    return res
  }

  *dp(): Generator<[number, number], void, unknown> {
    for (const v of this.bfsOrder) {
      if (v !== 0) {
        yield [this.suffixLink[v], v]
      }
    }
  }

  get size(): number {
    return this.children.length
  }
}
