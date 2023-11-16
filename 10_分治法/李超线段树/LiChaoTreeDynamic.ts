/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

type Line = {
  k: number
  b: number
}

type LiChaoNode = {
  lineId: number
  left: LiChaoNode | undefined
  right: LiChaoNode | undefined
}

const INF = 2e9 // 2e15

/**
 * 动态开点李超线段树.注意`添加线段`时空间消耗较大.
 */
class LiChaoTreeDynamic {
  private readonly _start: number
  private readonly _end: number
  private readonly _minimize: boolean
  private readonly _evaluate: (line: Line, x: number) => number
  private readonly _lines: Line[] = []
  private readonly _root = LiChaoTreeDynamic._initNode()

  constructor(
    start: number,
    end: number,
    options: {
      minimize?: boolean
      evaluate?: (line: Line, x: number) => number
    } = {}
  ) {
    const { minimize = true, evaluate = (line: Line, x: number) => line.k * x + line.b } = options
    this._start = start
    this._end = end
    this._minimize = minimize
    this._evaluate = evaluate
  }

  /** O(logn) */
  addLine(line: Line): void {
    const id = this._lines.length
    this._lines.push(line)
    this._addLine(this._root, id, this._start, this._end)
  }

  /** [start, end). O(log^2n) */
  addSegment(startX: number, endX: number, line: Line): void {
    if (startX >= endX) return
    const id = this._lines.length
    this._lines.push(line)
    this._addSegment(this._root, startX, endX, id, this._start, this._end)
  }

  /** O(logn) */
  query(x: number): { value: number; lineId: number } {
    if (!(this._start <= x && x < this._end)) {
      throw new RangeError(`x is out of range : ${x}`)
    }
    return this._query(this._root, x, this._start, this._end)
  }

  clear(): void {
    this._lines.length = 0
  }

  private _evaluateInner(fid: number, x: number): number {
    if (fid == -1) {
      return this._minimize ? INF : -INF
    }
    return this._evaluate(this._lines[fid], x)
  }

  private _addLine(node: LiChaoNode, fid: number, nodeL: number, nodeR: number): LiChaoNode {
    const gid = node.lineId
    const fl = this._evaluateInner(fid, nodeL)
    const fr = this._evaluateInner(fid, nodeR - 1)
    const gl = this._evaluateInner(gid, nodeL)
    const gr = this._evaluateInner(gid, nodeR - 1)
    const bl = this._minimize ? fl < gl : fl > gl
    const br = this._minimize ? fr < gr : fr > gr
    if (bl && br) {
      node.lineId = fid
      return node
    }
    if (!bl && !br) {
      return node
    }

    const nodeM = Math.floor((nodeL + nodeR) / 2)
    const fm = this._evaluateInner(fid, nodeM)
    const gm = this._evaluateInner(gid, nodeM)
    const bm = this._minimize ? fm < gm : fm > gm
    if (bm) {
      node.lineId = fid
      if (bl) {
        if (node.right == undefined) {
          node.right = LiChaoTreeDynamic._initNode()
        }
        node.right = this._addLine(node.right, gid, nodeM, nodeR)
      } else {
        if (node.left == undefined) {
          node.left = LiChaoTreeDynamic._initNode()
        }
        node.left = this._addLine(node.left, gid, nodeL, nodeM)
      }
    } else if (!bl) {
      if (node.right == undefined) {
        node.right = LiChaoTreeDynamic._initNode()
      }
      node.right = this._addLine(node.right, fid, nodeM, nodeR)
    } else {
      if (node.left == undefined) {
        node.left = LiChaoTreeDynamic._initNode()
      }
      node.left = this._addLine(node.left, fid, nodeL, nodeM)
    }
    return node
  }

  private _addSegment(node: LiChaoNode, xl: number, xr: number, fid: number, nodeL: number, nodeR: number): LiChaoNode {
    if (nodeL > xl) xl = nodeL
    if (nodeR < xr) xr = nodeR
    if (xl >= xr) return node
    if (nodeL < xl || xr < nodeR) {
      const nodeM = Math.floor((nodeL + nodeR) / 2)
      if (node.left == undefined) {
        node.left = LiChaoTreeDynamic._initNode()
      }
      if (node.right == undefined) {
        node.right = LiChaoTreeDynamic._initNode()
      }
      node.left = this._addSegment(node.left, xl, xr, fid, nodeL, nodeM)
      node.right = this._addSegment(node.right, xl, xr, fid, nodeM, nodeR)
      return node
    }
    return this._addLine(node, fid, nodeL, nodeR)
  }

  private _query(node: LiChaoNode, x: number, nodeL: number, nodeR: number): { value: number; lineId: number } {
    const fid = node.lineId
    let res = { value: this._evaluateInner(fid, x), lineId: fid }
    const nodeM = Math.floor((nodeL + nodeR) / 2)
    if (x < nodeM && node.left != undefined) {
      const cand = this._query(node.left, x, nodeL, nodeM)
      if (this._minimize ? cand.value < res.value : cand.value > res.value) {
        res = cand
      }
    }
    if (x >= nodeM && node.right != undefined) {
      const cand = this._query(node.right, x, nodeM, nodeR)
      if (this._minimize ? cand.value < res.value : cand.value > res.value) {
        res = cand
      }
    }
    return res
  }

  private static _initNode(): LiChaoNode {
    return { lineId: -1, left: undefined, right: undefined }
  }
}

export { LiChaoTreeDynamic }

if (require.main === module) {
  checkWithBf()

  function checkWithBf(): void {
    class Mocker {
      private readonly _minimize: boolean
      private readonly _lines: { start: number; end: number; line: Line }[] = []

      constructor(minimize: boolean) {
        this._minimize = minimize
      }

      addLine(line: Line): void {
        this._lines.push({ start: -Infinity, end: Infinity, line })
      }

      addSegment(start: number, end: number, line: Line): void {
        this._lines.push({ start, end, line })
      }

      query(x: number): { value: number; lineId: number } {
        let resValue = this._minimize ? INF : -INF
        let resLineId = -1
        for (const [index, { start, end, line }] of this._lines.entries()) {
          if (x >= start && x < end) {
            const value = line.k * x + line.b
            if (this._minimize ? value < resValue : value > resValue) {
              resValue = value
              resLineId = index
            }
          }
        }
        return { value: resValue, lineId: resLineId }
      }
    }

    const q = 500
    const points = Array(q)
      .fill(0)
      .map(() => Math.floor(-Math.random() * 1e5) + 5e4)
    const tree1 = new LiChaoTreeDynamic(-1e5, 1e5, { minimize: true })
    const tree2 = new Mocker(true)
    for (let i = 0; i < q; i++) {
      const k = -Math.floor(Math.random() * 1e5) + 5e4
      const b = -Math.floor(Math.random() * 1e5) + 5e4
      tree1.addLine({ k, b })
      tree2.addLine({ k, b })
      if (Math.random() < 0.5) {
        const k = -Math.floor(Math.random() * 1e5) + 5e4
        const b = -Math.floor(Math.random() * 1e5) + 5e4
        const start = Math.floor(Math.random() * 1e5)
        const end = Math.floor(Math.random() * 1e5)
        tree1.addSegment(start, end, { k, b })
        tree2.addSegment(start, end, { k, b })
      }

      const x = points[i]
      if (tree1.query(x).value !== tree2.query(x).value) {
        throw new Error()
      }
    }

    console.log('pass!')
  }
}
