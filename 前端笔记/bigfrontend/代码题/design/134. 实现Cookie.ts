declare interface Document {
  // "delete" 运算符的操作数必须是可选的。
  myCookie?: PropertyDescriptor
}

interface Cookie {
  value: string
  expires: number
}

// enable myCookie
function install() {
  const store = new Map<string, Cookie>()

  Object.defineProperty(document, 'myCookie', {
    configurable: true,
    get() {
      const cookies: string[] = []
      const time = Date.now()

      for (const [name, { value, expires }] of store) {
        if (expires <= time) store.delete(name)
        else cookies.push(`${name}=${value}`)
      }

      return cookies.join('; ')
    },
    set(val: string) {
      const { name, value, options } = parseCookie(val)
      let expires = Infinity
      if (options['max-age']) expires = Date.now() + Number(options['max-age']) * 1000
      store.set(name, { value, expires })
    },
  })
}

// disable myCookie
function uninstall() {
  delete document.myCookie
}

// name=value;max-age=30;secure=true;path='foo'
function parseCookie(str: string): {
  name: string
  value: string
  options: Record<string, unknown>
} {
  const [nameAndValue, ...rest] = str.split(';')
  const [name, value] = parseEntry(nameAndValue)

  const options = {} as Record<string, unknown>
  for (const option of rest) {
    const [key, value] = parseEntry(option)
    options[key] = value
  }

  return {
    name,
    value,
    options,
  }
}

function parseEntry(entry: `${string}=${string}`): [string, string]
function parseEntry(entry: string): string[]
function parseEntry(str: string): string[] {
  return str.split('=').map(s => s.trim())
}

// 不应该使用setTimeout(浏览器关闭时就没了，不能完全清除) 清除cookie 而是在 get时清除cookie
// image when you set a cookie, and use a setTimeout to clear it after 3 seconds, but you close your browser 1 second after now
// setTimeout is terminated, not guaranteed to finish
// when get is called, it yields the un-deleted bad data.
// That's why setTimeout should not be used. it is not accurate at all considering Event Loop
