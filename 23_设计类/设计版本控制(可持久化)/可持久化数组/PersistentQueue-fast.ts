/* eslint-disable no-param-reassign */

/**
 * 可持久化队列.
 * 初始时版本为0.
 */
class PersistentQueue<E = number> {
  private readonly _log: number
  private readonly _data: E[]
  private readonly _par: Uint32Array[]
  private readonly _backId: Uint32Array
  private readonly _size: Uint32Array
  private _curVersion = 0
  private _dataPtr = 0
  private _parPtr = 0
  private _backIdPtr = 1
  private _sizePtr = 1

  constructor(maxMutation: number) {
    maxMutation++
    this._log = 32 - Math.clz32(maxMutation)
    this._data = Array(maxMutation).fill(0)
    this._par = Array(maxMutation).fill(0)
    this._backId = new Uint32Array(maxMutation)
    this._size = new Uint32Array(maxMutation)
  }

  push(version: number, value: E): number {
    this._checkVersion(version)
    this._curVersion++
    const newId = this._dataPtr
    this._data[this._dataPtr++] = value
    this._par[this._parPtr++] = new Uint32Array(this._log)
    this._par[newId][0] = this._backId[version]
    this._backId[this._backIdPtr++] = newId
    this._size[this._sizePtr++] = this._size[version] + 1
    for (let d = 1; d < this._log; d++) {
      this._par[newId][d] = this._par[this._par[newId][d - 1]][d - 1]
    }
    return this._curVersion
  }

  shift(version: number): [newVersion: number, res: E | undefined] {
    this._curVersion++
    let r = this._backId[version]
    const len = this._size[version] - 1
    this._backId[this._backIdPtr++] = r
    this._size[this._sizePtr++] = len
    for (let d = 0; d < this._log; d++) {
      if ((len >> d) & 1) {
        r = this._par[r][d]
      }
    }
    return [this._curVersion, this._data[r]]
  }

  front(version: number): E | undefined {
    this._checkVersion(version)
    let r = this._backId[version]
    const len = this._size[version] - 1
    for (let d = 0; d < this._log; d++) {
      if ((len >> d) & 1) {
        r = this._par[r][d]
      }
    }
    return this._data[r]
  }

  back(version: number): E | undefined {
    this._checkVersion(version)
    return this._data[this._backId[version]]
  }

  private _checkVersion(version: number): void {
    if (version < 0) throw new Error('version must be non-negative')
    if (version > this._curVersion) {
      throw new Error('version must be less than or equal to current version')
    }
  }

  get curVersion(): number {
    return this._curVersion
  }
}

export { PersistentQueue }

if (require.main === module) {
  const queue = new PersistentQueue(50)
  let curV = 0
  curV = queue.push(curV, 10)
  curV = queue.push(curV, 100)
  curV = queue.push(curV, 200)
  curV = queue.push(curV, 300)
  console.log(queue.front(curV), queue.back(curV))

  const [curV2] = queue.shift(curV)
  console.log(queue.front(curV2), queue.back(curV2))
  const [curV3] = queue.shift(curV2)
  console.log(queue.front(curV3), queue.back(curV3))
}
