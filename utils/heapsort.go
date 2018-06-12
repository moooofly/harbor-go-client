package utils

import (
	"container/heap"
	"fmt"
)

type tagItem struct {
	tagName   string
	timestamp int64
}

type maxheap []*tagItem

func (h maxheap) Len() int { return len(h) }

func (h maxheap) Less(i, j int) bool {
	return h[i].timestamp > h[j].timestamp
}

func (h maxheap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *maxheap) Push(x interface{}) {
	it := x.(*tagItem)
	*h = append(*h, it)
}

func (h *maxheap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[0 : n-1]
	return it
}

// Sort on UNIX timestamp of tags
var maxh maxheap

type repoItem struct {
	data  *repoTop
	score float32
}

type minheap []*repoItem

func (h minheap) Len() int { return len(h) }

func (h minheap) Less(i, j int) bool {
	return h[i].score < h[j].score
}

func (h minheap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *minheap) Push(x interface{}) {
	it := x.(*repoItem)
	*h = append(*h, it)
}

func (h *minheap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[0 : n-1]
	return it
}

// Sort on the score of repos
var minh = minheap{}
var mhBk = minheap{}

func example1() {

	h := minheap{
		&repoItem{
			data: &repoTop{
				ID:           31,
				Name:         "prj2/photon",
				ProjectID:    3,
				Description:  "",
				PullCount:    5,
				StarCount:    0,
				TagsCount:    3,
				CreationTime: "2017-11-08T07:28:23Z",
				UpdateTime:   "2017-11-08T07:31:08Z",
			},
			score: 0.830000,
		},
		&repoItem{
			data: &repoTop{
				ID:           30,
				Name:         "sf3prj/hello-world",
				ProjectID:    13,
				Description:  "",
				PullCount:    0,
				StarCount:    0,
				TagsCount:    3,
				CreationTime: "2017-11-02T09:22:31Z",
				UpdateTime:   "2017-11-02T09:22:31Z",
			},
			score: 0.780000,
		},
		&repoItem{
			data: &repoTop{
				ID:           29,
				Name:         "prj5/hello-world",
				ProjectID:    6,
				Description:  "",
				PullCount:    4,
				StarCount:    0,
				TagsCount:    4,
				CreationTime: "2017-11-02T09:17:53Z",
				UpdateTime:   "2017-11-08T07:32:11Z",
			},
			score: 0.830000,
		},
		&repoItem{
			data: &repoTop{
				ID:           4,
				Name:         "prj3/hello-world",
				ProjectID:    4,
				Description:  "",
				PullCount:    4,
				StarCount:    0,
				TagsCount:    3,
				CreationTime: "2017-10-30T03:53:03Z",
				UpdateTime:   "2017-11-08T07:32:23Z",
			},
			score: 0.830000,
		},
	}

	heap.Init(&h)

	for h.Len() > 0 {
		item := heap.Pop(&h).(*repoItem)
		fmt.Printf("%.2f <==> %v\n", item.score, item.data)
	}
}
