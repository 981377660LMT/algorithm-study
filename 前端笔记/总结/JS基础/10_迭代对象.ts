const iterableObj = {
  a: 1,
  b: 2,
  c: 3,
  [Symbol.iterator]() {
    const entries = Object.entries(this)
    return {
      index: 0,
      next() {
        return {
          value: entries[this.index],
          done: this.index++ >= entries.length,
        }
      },
    }
  },
}

for (const item of iterableObj) {
  console.log(item)
}
