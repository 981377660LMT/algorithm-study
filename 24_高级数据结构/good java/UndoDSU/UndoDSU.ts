interface IOperation {
  apply(): void
  undo(): void
}

class UndoDSU {
  private readonly _rank: Int32Array
  private readonly _parent: Int32Array

  constructor(n: number) {
    this._rank = new Int32Array(n)
    this._parent = new Int32Array(n)
    for (let i = 0; i < n; i++) {
      this._rank[i] = 1
      this._parent[i] = -1
    }
  }

  init(): void {
    this._rank.fill(1)
    this._parent.fill(-1)
  }

  find(x: number): number {
    while (this._parent[x] !== -1) x = this._parent[x]
    return x
  }

  size(x: number): number {
    return this._rank[this.find(x)]
  }

  union(a: number, b: number): IOperation {
    let x: number
    let y: number

    return {
      apply: () => {
        x = this.find(a)
        y = this.find(b)
        if (x === y) return
        if (this._rank[x] < this._rank[y]) {
          const tmp = x
          x = y
          y = tmp
        }
        this._parent[y] = x
        this._rank[x] += this._rank[y]
      },
      undo: () => {
        let cur = y
        while (this._parent[cur] !== -1) {
          cur = this._parent[cur]
          this._rank[cur] -= this._rank[y]
        }
        this._parent[y] = -1
      }
    }
  }
}

export {}

if (require.main === module) {
  const n = 10
  const dsu = new UndoDSU(n)
  const operation = dsu.union(1, 2)
  operation.apply()
  console.log(dsu.size(1) === 2)
  operation.undo()
  console.log(dsu.size(1) === 1)
  console.log(dsu.size(2) === 1)
}
