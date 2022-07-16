class Set<T extends PropertyKey = any> {
  size: number
  private counter: Record<T, boolean>

  constructor(values?: readonly T[]) {
    this.size = 0
    this.counter = Object.create(null)
    for (const value of values ?? []) {
      this.add(value)
    }
  }

  has(value: T): boolean {
    return this.counter.hasOwnProperty(value)
  }

  add(value: T): this {
    if (this.has(value)) return this
    this.size++
    this.counter[value] = true
    return this
  }

  delete(value: T): boolean {
    if (!this.has(value)) return false
    this.size--
    delete this.counter[value]
    return true
  }
}

export {}
