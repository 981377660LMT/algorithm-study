/**
 *  One of the differences between null and undefined is how they are
    treated differently in JSON.stringify().

    ```js
    // JSON.stringify({a: null})      // '{"a":null}'
    // JSON.stringify({a: undefined}) // '{}'

    // JSON.stringify([null])         // '[null]'
    // JSON.stringify([undefined])    // '[null]'
    ```

    This difference might create troubles if there are missing alignments
    between client and server. It might be helpful to enforce using only one of them.
    You are asked to implement undefinedToNull() to return a copy
    that has all undefined replaced with null.

    ```js
    // undefinedToNull({a: undefined, b: 'BFE.dev'})
    // {a: null, b: 'BFE.dev'}

    // undefinedToNull({a: ['BFE.dev', undefined, 'bigfrontend.dev']})
    // {a: ['BFE.dev', null, 'bigfrontend.dev']}
    ```
 * https://bigfrontend.dev/zh/problem/undefined-to-null
 * !将对象中所有的 undefined 转换为 null
 */
function undefinedToNull(arg: unknown): unknown {
  dfs(arg)
  return arg

  // traverse the object and replace undefined with null
  function dfs(cur: unknown): void {
    if (typeof cur !== 'object' || cur === null) return
    for (const key of Object.keys(cur)) {
      if (cur[key] === undefined) {
        cur[key] = null
      } else {
        dfs(cur[key])
      }
    }
  }
}

// for (const next of Object.keys(cur)) 可以遍历Record对象的key
