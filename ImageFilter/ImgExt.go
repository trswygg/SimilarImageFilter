package imagefilter

import (
	"SimilarImageFilter/ImageFilter/Library/goimagehash"
	"SimilarImageFilter/ImageFilter/fileengine"
	"container/list"
	"fmt"
	"github.com/schollz/progressbar/v3"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	UseFileengine = 1
	UseveryThing  = 2
)

// 保存具体某一类hash
type hashExt struct {
	FileName string
	Hash     goimagehash.ExtImageHash
}

type SimilarPicture struct {
	A  string
	B  []im
	By byte
}
type im struct {
	name string
	dist int
}

//// trim repeated path
//func (group SimilarPicture) trim()  {
//	dirA,_ := filepath.Split(group.A)
//	for i,im:= range group.B {
//		dirB,ele := filepath.Split(im.name)
//		if dirA == dirB {
//			group.B[i].name = "~/" + ele
//		}
//	}
//}

// Comparator 图像对比引擎
type Comparator struct {
	// save the imageObj
	images []string
	count  func() int

	aHashList []hashExt
	dHashList []hashExt
	pHashList []hashExt

	config conf

	similarPictures *list.List
}

type conf struct {
	width  int
	height int
	dist   int

	// 1:fileEngine 2:EveryThing
	sourceEngine int
	// thread number default:4
	thread uint8

	// max image count
	max int

	ahash bool
	dhash bool
	phash bool
}

type Option func(p *Comparator)

func SetWidth(width int) Option {
	return func(p *Comparator) {
		p.config.width = width
	}
}
func SetHeight(height int) Option {
	return func(p *Comparator) {
		p.config.height = height
	}
}
func SetDist(distance int) Option {
	return func(p *Comparator) {
		p.config.dist = distance
	}
}
func SetSource(source int) Option {
	return func(p *Comparator) {
		p.config.sourceEngine = source
	}
}
func SetThread(threadNum int) Option {
	return func(p *Comparator) {
		p.config.thread = uint8(threadNum)
	}
}
func SetMax(max int) Option {
	return func(p *Comparator) {
		p.config.max = max
	}
}
func SetHash(hash string) Option {
	hash = strings.TrimSpace(strings.ToLower(hash))
	ahash, dhash, phash := false, false, false
	for _, op := range hash {
		switch op {
		case 'a':
			ahash = true
		case 'd':
			dhash = true
		case 'p':
			phash = true
		}
	}
	return func(p *Comparator) {
		p.config.ahash = ahash
		p.config.dhash = dhash
		p.config.phash = phash
	}
}

//NewImageFilter  return a new Comparator
func NewImageFilter(options ...Option) *Comparator {
	c := Comparator{
		images: make([]string, 0, 2000),
		config: conf{
			width:        64,
			height:       64,
			dist:         128,
			sourceEngine: 1,
			max:          32767, // max of int32
			ahash:        true,
			dhash:        true,
			phash:        true,
		},
	}
	c.count = func() int {
		return len(c.images)
	}
	for _, op := range options {
		op(&c)
	}
	// init
	c.aHashList = make([]hashExt, 0, 2000)
	c.dHashList = make([]hashExt, 0, 2000)
	c.pHashList = make([]hashExt, 0, 2000)
	c.similarPictures = list.New()
	return &c
}

// toProcess schedule  fill images
func (C *Comparator) toProcess(WorkerCount int) {
	bar := basicBar(C.count(), "[cyan][][reset] Processing...")
	//defer bar.Finish()
	wg := sync.WaitGroup{}
	tasks := splitArray(C.images, WorkerCount)
	fmt.Printf("[INFO] Start process images:%d \n", WorkerCount)
	fmt.Printf("[DEBUG] splitArray C.images [%d] -> [][]string len:%d, cap:%d \n",
		len(C.images), len(tasks), cap(tasks))
	for i := 0; i < len(tasks); i++ {
		wg.Add(1)
		go C.processExt(tasks[i][:], &wg, bar)
	}
	wg.Wait()
	fmt.Print("\n")
}

// goroutines processExt
func (C *Comparator) processExt(task []string, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()

	w, h := C.config.width, C.config.height
	for _, path := range task {
		_ = bar.Add(1)
		i, err := decodeImg(path)
		if err != nil || i == nil {
			// pass
			continue
		}

		// process
		if C.config.ahash == true {
			ahash, err := goimagehash.ExtAverageHash(i, w, h)
			if err != nil {
				fmt.Print(err.Error())
				continue
			}
			h := hashExt{
				FileName: path,
				Hash:     *ahash,
			}
			C.aHashList = append(C.aHashList, h)
		}
		if C.config.dhash == true {
			dhash, err := goimagehash.ExtDifferenceHash(i, w, h)
			if err != nil {
				fmt.Print(err.Error())
				continue
			}
			h := hashExt{
				FileName: path,
				Hash:     *dhash,
			}
			C.dHashList = append(C.dHashList, h)
		}
		if C.config.phash == true {
			phash, err := goimagehash.ExtPerceptionHash(i, w, h)
			if err != nil {
				fmt.Print(err.Error())
				continue
			}
			h := hashExt{
				FileName: path,
				Hash:     *phash,
			}
			C.pHashList = append(C.pHashList, h)
		}
	}
}

// Index Try to index all image files from ‘path’
func (C *Comparator) Index(path string) {
	fmt.Printf("[INFO] Start index on path:%s \n", path)
	if isDir(path) == true {
		switch C.config.sourceEngine {
		case UseFileengine:
			C.indexFromFileEngine(path)
		case UseveryThing:
			C.indexFromEveryThing(path)
		default:
			C.indexFromFileEngine(path)
		}
		fmt.Printf("[INFO] Index %d Objects \n", C.count())
	} else {
		fmt.Printf("[ERROR] path:%s is not a directory \n", path)
	}
	if C.count() > C.config.max {
		C.images = C.images[:C.config.max]
		fmt.Printf("[INFO] Objects overrang，resize %d -> %d \n", C.count(), len(C.images))
	}
	C.toProcess(runtime.NumCPU())
}

func (C *Comparator) indexFromFileEngine(path string) {
	engine := fileengine.DefaultImgEngine(&C.images)
	res, count := engine.Start(path)
	C.images = *res
	if count != C.count() {
		fmt.Printf("[ERROR] indexFromFileEngine C.count: %d, count: %d", C.count(), count)
	}
}

// WARNING: Unrealized
func (C *Comparator) indexFromEveryThing(path string) {
	engine := fileengine.DefaultImgEngine(&C.images)
	res, count := engine.Start(path)
	C.images = *res
	if count != C.count() {
		fmt.Printf("[ERROR] indexFromEveryThing C.count: %d, count: %d", C.count(), count)
	}
}

// 遍历整个列表需找汉明距离小于`dist` 的一组图片
// test: speed count?
func (C *Comparator) Contras() {
	wg := sync.WaitGroup{}
	dist := C.config.dist
	if C.config.ahash {
		wg.Add(1)
		go func() {
			C.contrasByAHash(dist)
			defer wg.Done()
		}()
	}
	if C.config.dhash {
		wg.Add(1)
		go func() {
			C.contrasByDHash(dist)
			defer wg.Done()
		}()
	}
	if C.config.phash {
		wg.Add(1)
		go func() {
			C.contrasByPHash(dist)
			defer wg.Done()
		}()
	}
	wg.Wait()
	// delete unused HashList
	C.aHashList = nil
	C.dHashList = nil
	C.pHashList = nil
}
func (C *Comparator) contrasByAHash(dist int) {
	for i := 0; i < len(C.aHashList); i++ {
		// should lazy malloc?
		group := SimilarPicture{
			A:  C.aHashList[i].FileName,
			B:  nil,
			By: 'a',
		}

		for j := 0; j < len(C.aHashList); j++ {
			if i == j {
				continue
			}
			d, _ := C.aHashList[i].Hash.Distance(&C.aHashList[j].Hash)
			if d < dist {
				group.B = append(group.B, im{
					name: C.aHashList[j].FileName,
					dist: d,
				})
			}
		}
		if group.B != nil {
			C.similarPictures.PushBack(group)
		}
	}
}
func (C *Comparator) contrasByDHash(dist int) {
	for i := 0; i < len(C.dHashList); i++ {
		// should lazy malloc?
		group := SimilarPicture{
			A:  C.dHashList[i].FileName,
			B:  nil,
			By: 'd',
		}

		for j := 0; j < len(C.dHashList); j++ {
			if i == j {
				continue
			}
			d, _ := C.dHashList[i].Hash.Distance(&C.dHashList[j].Hash)
			if d < dist {
				group.B = append(group.B, im{
					name: C.dHashList[j].FileName,
					dist: d,
				})
			}
		}
		if group.B != nil {
			C.similarPictures.PushBack(group)
		}
	}
}
func (C *Comparator) contrasByPHash(dist int) {
	for i := 0; i < len(C.dHashList); i++ {
		// should lazy malloc?
		group := SimilarPicture{
			A:  C.pHashList[i].FileName,
			B:  nil,
			By: 'p',
		}

		for j := 0; j < len(C.pHashList); j++ {
			if i == j {
				continue
			}
			d, _ := C.pHashList[i].Hash.Distance(&C.pHashList[j].Hash)
			if d < dist {
				group.B = append(group.B, im{
					name: C.pHashList[j].FileName,
					dist: d,
				})
			}
		}
		if group.B != nil {
			C.similarPictures.PushBack(group)
		}
	}
}

func (C *Comparator) trimResult(by byte) {
	// only run when corresponding hash is enable
	if !C.config.ahash && by == 'a' || !C.config.dhash && by == 'd' || !C.config.phash && by == 'p' {
		return
	}
	// save result
	m := make(map[string]string, 10)
	for group := C.similarPictures.Front(); group != nil; group = group.Next() {
		g := group.Value.(SimilarPicture)
		if g.By != by {
			continue
		}
		for _, im := range g.B {
			m[im.name] = g.A
		}
	}
	// trim repeated comp
	var n *list.Element
	for group := C.similarPictures.Front(); group != nil; group = n {
		g := group.Value.(SimilarPicture)
		n = group.Next()
		if g.By != by {
			continue
		}
		if len(g.B) == 1 {
			value, ok := m[g.A]
			if ok && value == g.B[0].name {
				//hit
				//fmt.Printf("should remove %s %v By:%c \n", g.A, g.B, g.By)
				C.similarPictures.Remove(group)
				delete(m, g.B[0].name)
			}
		}
		// trim
		dirA, _ := filepath.Split(g.A)
		for i, im := range g.B {
			dirB, ele := filepath.Split(im.name)
			if dirA == dirB {
				g.B[i].name = "./" + ele
			}
		}
	}
}