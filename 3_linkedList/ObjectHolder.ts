class ObjectHolder<T> {
  value: T | undefined

  clear(): void {
    this.value = undefined
  }
}

export { ObjectHolder }
