type Obj = Record<PropertyKey, unknown> | unknown[]

const arrayHandler: ProxyHandler<unknown[]> = {
  set(_, prop) {
    console.log(prop)
    return true
  }
}

const objectHandler: ProxyHandler<Record<PropertyKey, any>> = {
  set(_, prop) {
    console.log(prop)
    return true
  }
}

const propHandler: ProxyHandler<Obj> = {
  get(target, p) {
    const val = Reflect.get(target, p)
    if (!isObj(val)) return val
    return lazyProxy(val)
  },
  set(target, prop) {
    if (Array.isArray(target)) {
      console.log(prop)
      return true
    }
    console.log(prop)
    return true
  }
}

const methodHandler: ProxyHandler<(...args: unknown[]) => unknown> = {
  apply(target) {
    console.log(target.name)
  }
}

function isObj(obj: unknown): obj is Obj {
  return typeof obj === 'object' && obj !== null
}

/**
 * 初始化时，遍历对象的所有属性，将其转换为代理对象.
 */
function deepProxy<T extends Obj>(obj: T): T {
  if (Array.isArray(obj)) {
    obj.forEach((val, index) => {
      if (isObj(val)) {
        obj[index] = deepProxy(val)
      }
    })
    return new Proxy(obj, arrayHandler) as T
  }

  const keys = Object.keys(obj)
  keys.forEach(key => {
    const val = obj[key]
    if (isObj(val)) {
      obj[key] = deepProxy(val)
    }
  })
  return new Proxy(obj, objectHandler)
}

/**
 * 惰性代理，只有在访问属性时，才会将其转换为代理对象.
 */
function lazyProxy<T extends Obj>(obj: T): T {
  if (Array.isArray(obj)) {
    console.log('array')
  }
  return new Proxy(obj, propHandler) as T
}

export {}
