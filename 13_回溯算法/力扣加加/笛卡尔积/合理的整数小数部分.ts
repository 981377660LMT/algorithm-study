const isValidInteger = (str: string) => str.length === 1 || str[0] !== '0'
const isValidDecimal = (str: string) => str[str.length - 1] !== '0'

export { isValidDecimal, isValidInteger }
