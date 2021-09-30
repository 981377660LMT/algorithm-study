class DefaultDict<K extends PropertyKey, V = unknown> extends Object {
  private defaultFactory: (...args: any[]) => V

  constructor(defaultFactory: (...args: any[]) => V, objectValue?: any) {
    super(objectValue)
    this.defaultFactory = defaultFactory
    // @ts-ignore
    return this.buildDict()
  }

  private buildDict(): Record<K, V> {
    const dict = Object.create(null)

    return new Proxy(dict, {
      get: (target, key, receiver) => {
        if (!Reflect.has(target, key)) {
          Reflect.set(target, key, this.defaultFactory(key))
        }

        return Reflect.get(target, key)
      },
    })
  }
}

if (require.main === module) {
  const dict = new DefaultDict(Number)
  // @ts-ignore
  console.log(dict[1][4])
  console.log(dict)
}

export { DefaultDict }
