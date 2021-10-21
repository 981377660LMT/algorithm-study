import { Crawler } from '../Crawler'

class YangCrawler extends Crawler {
  protected override async normalizeData(data: string[]): Promise<string[]> {
    return data
  }
}

const yangCrawler = new YangCrawler({
  url: 'https://muyiy.cn/question/',
  xPathExpression: "//p[@class='sidebar-heading']//text()",
  dirName: __dirname,
})

yangCrawler.crawl().then(console.log).catch(console.log)
