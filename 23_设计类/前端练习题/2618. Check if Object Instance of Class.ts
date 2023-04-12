function checkIfInstanceOf(obj: any, classFunction: any): boolean {
  if (obj == null || typeof classFunction !== 'function') return false

  while (obj != null) {
    // getPrototype of in place of __proto__
    const proto = Object.getPrototypeOf(obj)
    if (proto === classFunction.prototype) return true
    obj = proto
  }

  return false
}

export {}
