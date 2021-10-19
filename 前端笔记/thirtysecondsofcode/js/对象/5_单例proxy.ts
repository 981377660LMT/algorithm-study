const singletonify = (className: Function) => {
  return new Proxy(className.prototype.constructor, {
    // @ts-ignore
    instance: null,
    // new 的时候屌用
    construct(target, argumentsList) {
      if (!this.instance) this.instance = new target(...argumentsList)
      return this.instance
    },
  })
}
