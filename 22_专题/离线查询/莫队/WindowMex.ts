/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { SortedListFast } from '../根号分治/SortedList/SortedListFast'

/**
 * @alias MexManager
 */
class WindowMex {
  private readonly _maxOperation: number
  private readonly _mexStart: number
  private readonly _counter: Uint32Array
  private readonly _sl: SortedListFast<number>

  constructor(maxOperation: number, mexStart = 0) {
    this._maxOperation = maxOperation
    this._mexStart = mexStart
    this._counter = new Uint32Array(maxOperation + 1)
    const initData = Array<number>(maxOperation + 1)
    for (let i = 0; i <= maxOperation; ++i) {
      initData[i] = mexStart + i
    }
    this._sl = new SortedListFast(initData)
  }

  add(v: number): boolean {
    if (v < this._mexStart || v > this._mexStart + this._maxOperation) {
      return false
    }
    this._counter[v - this._mexStart]++
    if (this._counter[v - this._mexStart] === 1) {
      this._sl.discard(v)
    }
    return true
  }

  discard(v: number): boolean {
    if (v < this._mexStart || v > this._mexStart + this._maxOperation) {
      return false
    }
    if (!this._counter[v - this._mexStart]) {
      return false
    }
    this._counter[v - this._mexStart]--
    if (!this._counter[v - this._mexStart]) {
      this._sl.add(v)
    }
    return true
  }

  query(): number {
    if (!this._sl.length) {
      return this._mexStart
    }
    return this._sl.min!
  }
}

class WindowMexNaive {
  private readonly _maxOperation: number
  private readonly _mexStart: number
  private readonly _counter: Uint32Array
  private _mex: number

  constructor(maxOperation: number, mexStart = 0) {
    this._maxOperation = maxOperation
    this._mexStart = mexStart
    this._counter = new Uint32Array(maxOperation + 1)
    this._mex = mexStart
  }

  add(v: number): boolean {
    if (v < this._mexStart || v > this._mexStart + this._maxOperation) {
      return false
    }
    this._counter[v - this._mexStart]++
    while (this._counter[this._mex - this._mexStart]) {
      this._mex++
    }
    return true
  }

  discard(v: number): boolean {
    if (v < this._mexStart || v > this._mexStart + this._maxOperation) {
      return false
    }
    if (!this._counter[v - this._mexStart]) {
      return false
    }
    this._counter[v - this._mexStart]--
    if (!this._counter[v - this._mexStart]) {
      if (v < this._mex) {
        this._mex = v
      }
    }
    return true
  }

  query(): number {
    return this._mex
  }
}

export { WindowMex }

if (require.main === module) {
  check()
  function check(): void {
    const MAX_OPERATION = 1e6
    const MEX_START = 0
    const mocker = new WindowMexNaive(MAX_OPERATION, MEX_START)
    const windowMex = new WindowMex(MAX_OPERATION, MEX_START)
    for (let i = 0; i < MAX_OPERATION; ++i) {
      const v = Math.floor(Math.random() * 1e9)
      const op = Math.random() < 0.8 ? 'add' : 'discard'
      mocker[op](v)
      windowMex[op](v)
      if (mocker.query() !== windowMex.query()) {
        console.log('error')
        console.log(mocker.query(), windowMex.query())
        console.log(mocker)
        console.log(windowMex)
        break
      }
    }
  }
  console.log('test passed')
}
