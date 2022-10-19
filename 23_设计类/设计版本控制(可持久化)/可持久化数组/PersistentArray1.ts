import assert from 'assert'

interface PersistentArray {
  /**
   * 当前最新的版本号 (从 0 开始, 0 代表初始数组)
   */
  readonly curVersion: number

  /**
   * @param version 版本号 0 <= version <= {@link curVersion}
   * @param index 数组下标 0 <= index < n
   * @returns 版本号为`version`的数组中下标为`index`的值
   */
  query: (version: number, index: number) => number

  /**
   * @param version 版本号 0 <= version <= {@link curVersion}
   * @param index 数组下标 0 <= index < n
   * @param value 更新的值
   * @returns 新数组的版本号
   */
  update: (version: number, index: number, value: number) => number
}

/**
 * 创建一个可持久化数组
 * @param sizeOrArray 数组大小或者数组
 * @param updateTimes 更新次数的上界
 */
function usePersistentArray(
  sizeOrArray: number | ArrayLike<number>,
  updateTimes: number
): PersistentArray {
  let _curVersion = 0
  let _nodeId = 0

  const isNumber = typeof sizeOrArray === 'number'
  const _n = isNumber ? sizeOrArray : sizeOrArray.length
  const _size = 4 * _n + Number(_n).toString(2).length * updateTimes

  const _leftChild = new Uint32Array(_size)
  const _rightChild = new Uint32Array(_size)
  const _treeValue = Array<number>(_size).fill(0)

  const _roots = new Uint32Array(updateTimes + 1)
  _roots[0] = _build(0, _n - 1, isNumber ? undefined : sizeOrArray)

  function query(version: number, index: number): number {
    return _query(_roots[version], 0, _n - 1, index)
  }

  function update(version: number, index: number, value: number): number {
    const rootId = _update(_roots[version], 0, _n - 1, index, value)
    _curVersion++
    _roots[_curVersion] = rootId
    return _curVersion
  }

  return {
    get curVersion() {
      return _curVersion
    },
    query,
    update
  }

  function _build(left: number, right: number, array?: ArrayLike<number>): number {
    const node = _nodeId
    _nodeId++
    if (left === right) {
      _treeValue[node] = array ? array[left] : 0
      return node
    }

    const mid = (left + right) >> 1
    _leftChild[node] = _build(left, mid, array)
    _rightChild[node] = _build(mid + 1, right, array)
    return node
  }

  function _query(curRoot: number, left: number, right: number, pos: number): number {
    if (left === right) {
      return _treeValue[curRoot]
    }

    const mid = (left + right) >> 1
    if (pos <= mid) {
      return _query(_leftChild[curRoot], left, mid, pos)
    }
    return _query(_rightChild[curRoot], mid + 1, right, pos)
  }

  function _update(
    preRoot: number,
    left: number,
    right: number,
    pos: number,
    value: number
  ): number {
    const node = _nodeId
    _nodeId++
    _leftChild[node] = _leftChild[preRoot]
    _rightChild[node] = _rightChild[preRoot]
    _treeValue[node] = _treeValue[preRoot]
    if (left === right) {
      _treeValue[node] = value
      return node
    }

    const mid = (left + right) >> 1
    if (pos <= mid) {
      _leftChild[node] = _update(_leftChild[preRoot], left, mid, pos, value)
    } else {
      _rightChild[node] = _update(_rightChild[preRoot], mid + 1, right, pos, value)
    }
    return node
  }
}

if (require.main === module) {
  const nums = [59, 46, 14, 87, 41]
  const v0 = 0
  const persistentArray = usePersistentArray(nums, 4)
  assert.strictEqual(persistentArray.query(v0, 0), 59)
  assert.strictEqual(persistentArray.query(v0, 1), 46)
  assert.strictEqual(persistentArray.query(v0, 2), 14)
  assert.strictEqual(persistentArray.query(v0, 3), 87)
  assert.strictEqual(persistentArray.query(v0, 4), 41)

  const v1 = persistentArray.update(v0, 0, 100)
  assert.strictEqual(persistentArray.query(v1, 0), 100)

  const v2 = persistentArray.update(v1, 1, 200)
  assert.strictEqual(persistentArray.query(v2, 1), 200)

  const v3 = persistentArray.update(v0, 2, 300)
  for (let i = 0; i < 5; i++) {
    console.log(persistentArray.query(v3, i))
  }
}

export { usePersistentArray, PersistentArray }
