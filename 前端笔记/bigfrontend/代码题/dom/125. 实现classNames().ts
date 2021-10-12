/**
 * @param {any[]} args
 * @returns {string}
 * classNames()接受任意的参数，其将过滤其中的falsy值，
 * 然后生成最终的class name字符串。
 */
function classNames(...args: any[]): string {
  // your code here
  return args
    .flat(Infinity)
    .reduce<(string | number)[]>((res, item) => {
      if (item === null) return res

      switch (typeof item) {
        case 'string':
        case 'number':
          res.push(item)
          break
        case 'object':
          for (const [key, value] of Object.entries(item)) {
            if (!!value) res.push(key)
          }
          break
        default:
          break
      }

      return res
    }, [])
    .join(' ')
}

// const classNames = require('classnames');

// class Button extends React.Component {
//   // ...
//   render () {
//     const btnClass = classNames({
//        btn: true,
//       'btn-pressed': this.state.isPressed,
//       'btn-over': !this.state.isPressed && this.state.isHovered
//     });
//     return <button className={btnClass}>{this.props.label}</button>;
//   }
// }

// 1.string 和 number 的话，直接使用
// classNames('BFE', 'dev', 100)
// // 'BFE dev 100'

// 2.其他的primitive将被忽略
// classNames(
//   null, undefined, Symbol(), 1n, true, false
// )
// // ''

// 3.Object的enumerable property，如果key是string而且value是truthy的话将被保留。数组需要扁平化。
const obj = new Map()
// @ts-ignore
obj.cool = '!'
console.log(Object.entries(obj)) // 竟然可以map.属性名来设置属性...
console.log(classNames({ BFE: [], dev: true, is: 3 }, obj))
// 'BFE dev is cool'
console.log(classNames(['BFE', [{ dev: true }, ['is', [obj]]]]))
// 'BFE dev is cool'

console.log(Object.entries([1, 2, 3])) // [ [ '0', 1 ], [ '1', 2 ], [ '2', 3 ] ]
console.log(Object.entries(new Map([[1, 2]]))) // []
