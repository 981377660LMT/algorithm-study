from collections import defaultdict
from typing import List


class Solution:
    def subdomainVisits(self, cpdomains: List[str]) -> List[str]:
        counter = defaultdict(int)
        for item in cpdomains:
            count, domain = item.split()
            count = int(count)
            splits = domain.split(".")
            for i in range(len(splits)):
                counter[".".join(splits[i:])] += count

        return ["{} {}".format(v, k) for k, v in counter.items()]


print(
    Solution().subdomainVisits(
        cpdomains=["900 google.mail.com", "50 yahoo.com", "1 intel.mail.com", "5 wiki.org"]
    )
)

# 输出：["901 mail.com","50 yahoo.com","900 google.mail.com","5 wiki.org","5 org","1 intel.mail.com","951 com"]
# 解释：按照前文描述，会访问 "google.mail.com" 900 次，"yahoo.com" 50 次，"intel.mail.com" 1 次，"wiki.org" 5 次。
# 而对于父域名，会访问 "mail.com" 900 + 1 = 901 次，"com" 900 + 50 + 1 = 951 次，和 "org" 5 次。
