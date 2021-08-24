package cache

// * 添加一个新node
// hashmap.put(node)
// 如果node的权重大于windowq的最大权重，push到windowq的first，否则push到windowq的last
// 如果window的当前权重大于window最大权重，挪动window的first，放到probation的last，直到window的当前权重小于等于window的最大权重。到此：window的当前权重已经收缩到合理值了。
// loop：如果cache的当前权重超出最大权重，进行淘汰：
//   如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。

// * 更新一个node
// hashmap.put(node)
// node.update Weight and value
// if node is belongs to windowq:
//   如果node的权重大于windowq的最大权重，移动到windowq的first，否则移动到windowq的last
//   如果window的当前权重大于window最大权重，挪动window的first，放到probation的last，直到window的当前权重小于等于window的最大权重。到此：window的当前权重已经收缩到合理值了。
// elif node is belongs to probationq:
//   挪动node到protected
//   如果protected的当前权重大于protected最大权重，挪动protected的first，放到probation的last，直到protected的当前权重小于等于protected的最大权重。到此：protected的当前权重已经收缩到合理值了。
//   loop：如果cache的当前权重超出最大权重，进行淘汰：
//     如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。
// elif node is belongs to protected:
//   挪动node到protected的last
//   如果protected的当前权重大于protected最大权重，挪动protected的first，放到probation的last，直到protected的当前权重小于等于protected的最大权重。到此：protected的当前权重已经收缩到合理值了。
//   loop：如果cache的当前权重超出最大权重，进行淘汰：
//     如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。

// * 获取一个key的value
// if hit in window:
//   挪动到window队尾，返回value
// elif hit in probation:
//   挪动node，从probation到protected的队尾
//   如果protected的当前权重大于protected最大权重，挪动protected的first，放到probation的last，直到protected的当前权重小于等于protected的最大权重。到此：protected的当前权重已经收缩到合理值了。
//   loop：如果cache的当前权重超出最大权重，进行淘汰：
//     如果probation的 victim(first) 和 candidate(last) 进行对比，按照FrequencyCandidate 和 FrequencyVictim 和 随机数 一起来判断淘汰 Victim 或者 Candidate。到此：Cache的当前权重已经收缩到合理值了。
// elif hit in protected:
//   挪动node到protected的队尾

// *删除一个node
// 直接从hashmap里面删除掉
// 从对应q里面删除掉
