type Product = number
type Price = number

class Cashier {
  private client: number
  private menu: Map<Product, Price>

  constructor(
    private n: number,
    private discount: number,
    private products: number[],
    private prices: number[]
  ) {
    this.client = 0
    this.menu = new Map()
    for (let i = 0; i < products.length; i++) {
      this.menu.set(products[i], prices[i])
    }
  }

  getBill(product: number[], amount: number[]): number {
    this.client++
    let res = 0
    for (let i = 0; i < product.length; i++) {
      res += this.menu.get(product[i])! * amount[i]
    }
    if (this.client % this.n === 0) {
      res -= (this.discount * res) / 100
    }
    return res
  }
}
