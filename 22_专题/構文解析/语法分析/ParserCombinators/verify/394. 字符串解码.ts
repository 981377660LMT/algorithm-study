import { lazy, oneOf, seqOf, str, zeroOrMore } from '../lib'

// encoded -> (Str | Num "[" encoded "]")*
const encoded = lazy(() => zeroOrMore(oneOf(str, seqOf())))
