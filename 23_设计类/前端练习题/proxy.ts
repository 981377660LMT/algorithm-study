/**
 * 初始化时，遍历对象的所有属性，将其转换为代理对象.
 */
function deepProxy<T>(obj: T): T {}

/**
 * 惰性代理，只有在访问属性时，才会将其转换为代理对象.
 */
function lazyProxy<T>(obj: T): T {}

export {}
