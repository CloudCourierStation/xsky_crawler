package campus_recruitment

import "testing"

func TestCrawler(t *testing.T) {
	err := Crawler("https://xskydata.jobs.feishu.cn/school",1,100)
	if err != nil {
		return
	}
}

func BenchmarkCrawler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := Crawler("https://xskydata.jobs.feishu.cn/school",1,100)
		if err != nil {
			b.Fatal(err)
			return
		}
	}
}
