class MyURLSearchParams implements URLSearchParams {
  private paramsMap: Map<string, string[]>

  constructor(init: string) {
    this.paramsMap = new Map()
    const searchParams = init.replace(/^\?/, '')
    const searchPairs = searchParams.split('&')

    searchPairs.forEach(pair => {
      const [key, value] = pair.split('=')
      this.append(key, value)
    })
  }

  append(name: string, value: string): void {
    !this.paramsMap.has(name) && this.paramsMap.set(name, [])
    this.paramsMap.get(name)!.push(`${value}`)
  }

  delete(name: string): void {
    this.paramsMap.delete(name)
  }

  /**
   * @param {string} name
   * returns the first value of the name
   */
  get(name: string): string | null {
    return this.paramsMap.get(name)?.[0] || null
  }

  /**
   * @param {string} name
   * @return {string[]}
   * returns the value list of the name
   */
  getAll(name: string): string[] {
    return this.paramsMap.get(name) || []
  }

  has(name: string): boolean {
    return this.paramsMap.has(name)
  }

  set(name: string, value: string): void {
    this.paramsMap.set(name, [`${value}`])
  }

  sort(): void {
    const sortedKeys = [...this.paramsMap.keys()].sort()
    const originalMap = this.paramsMap
    this.paramsMap = new Map()
    sortedKeys.forEach(k => this.paramsMap.set(k, originalMap.get(k)!))
  }

  toString(): string {
    return [...this.entries()].reduce(
      (pre, cur, index) => pre + (index !== 0 ? '&' : '') + cur.join('='),
      ''
    )
  }

  forEach(
    callbackfn: (value: string, key: string, parent: MyURLSearchParams) => void,
    thisArg?: any
  ): void {
    ;[...this.entries()].forEach(([k, v]) => callbackfn(v, k, this))
  }

  /**
   * @return {Iterator}
   */
  keys(): Iterator<string> {
    return this.paramsMap.keys()
  }

  /**
   * @return {Iterator} values
   */
  *values(): Generator<string, any, any> {
    for (const arr of this.paramsMap.values()) {
      yield* arr
    }
  }

  /**
   * @returns {Iterator}
   */
  *entries(): Generator<[string, string], any, any> {
    const result: [string, string][] = []
    for (const [key, arr] of this.paramsMap.entries()) {
      for (const value of arr) {
        result.push([key, value])
      }
    }
    yield* result
  }
}

const params = new MyURLSearchParams('?a=1&a=2&b=2')
console.log(params)
params.get('a') // '1'
params.getAll('a') // ['1', '2']
params.get('b') // '2'
params.getAll('b') // ['2']

params.append('a', '3')
params.set('b', '3')
console.log(params.toString()) // 'a=1&a=2&b=3&a=3'

// Iterable: 具有Symbol.iterator的 对象
// Iterator: 具有next/return/throw 的 对象
// Generator: 继承了Iterator 同时也是 Iterable 的，即IterableIterator
