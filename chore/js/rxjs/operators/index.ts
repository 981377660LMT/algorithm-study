// https://rxjs-cn.github.io/learn-rxjs-operators/

import { audit, buffer, fromEvent, interval, map } from 'rxjs'

const source = interval(1000)
const op = map((x: number) => x * x)
op(source).subscribe(console.log)
