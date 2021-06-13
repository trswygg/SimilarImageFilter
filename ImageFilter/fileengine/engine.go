package fileengine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// Icon :fileengine powered by trswygg
var Icon = `
    _______ __     ______            _
   / ____(_) /__  / ____/___  ______(_)___  ___
  / /_  / / / _ \/ __/ / __ \/ __  / / __ \/ _ \
 / __/ / / /  __/ /___/ / / / /_/ / / / / /  __/
/_/   /_/_/\___/_____/_/ /_/\__, /_/_/ /_/\___/
                           /____/
							powered by trswygg
`

// bar is a simple progress bar
var bar *progressbar.ProgressBar

// result is a slice save all  result
var r *[]string

// FileEngine 文件搜索引擎
type FileEngine struct {
	targetType []string
	ignorePath []string

	// engine core
	workerCount   int32
	maxWorker     int32
	searchRequest chan string
	foundMatch    chan string
	workerDone    chan bool
	m             sync.Mutex
}

// ------------------ Synchronized Iterator ------------------

// DefaultImgEngine 默认图片搜索引擎
func DefaultImgEngine(result *[]string) *FileEngine {
	engine := FileEngine{}
	engine.targetType = []string{".jpg", ".png"}
	engine.ignorePath = []string{"$RECYCLE.BIN", "Windows", "tmp", "Tmp", "System Volume Information"}

	engine.workerCount = 0
	engine.maxWorker = int32(4)
	engine.searchRequest = make(chan string)
	engine.workerDone = make(chan bool)
	engine.foundMatch = make(chan string, 16)

	r = result

	bar = basicBar(-1, "[light_blue]][1/8][reset] Searching...")

	return &engine
}

// Start :Start FileEngine on path
func (engine *FileEngine) Start(path string) (*[]string, int) {
	startTime := time.Now()
	fmt.Printf("[FileEngine] %s Start() INFO: Starting task on %s \n", startTime.Format("20060102 15:04:05"), path)

	// 修复 path 中的分界符
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	//go engine.search(path, true)
	engine.scheduler(path)

	fmt.Printf("\n[FileEngine] Start() End task ,TimeConsum %s, found %d \n", time.Since(startTime).String(), len(*r))
	return r, len(*r)
}

// scheduler 调度器
// 调度工作线程
func (engine *FileEngine) scheduler(path string) {
	maxWorker := engine.maxWorker
	//
	go engine.search(path, true)
	engine.workerCount++
	fmt.Printf("[FileEngine] scheduler() INFO: schedulering ,maxWorker = %d \n", maxWorker)
	bar := basicBar(-1, "[light_blue]][1/8][reset] Searching...")
	defer func() {
		_ = bar.Finish()
	}()

	// AnInfinite Loop
	for {
		select {

		// new task
		case path := <-engine.searchRequest:
			engine.m.Lock()
			engine.workerCount++
			engine.m.Unlock()
			go engine.search(path, true)
			//fmt.Printf("[FileEngine] scheduler() DEBUG: new search task on path:%s,left:%d/%d \n", path,maxWorker , engine.workerCount)

		// task done
		case <-engine.workerDone:
			engine.m.Lock()
			engine.workerCount--
			engine.m.Unlock()
			// fmt.Printf("[FileEngine] scheduler() DEBUG: search goroutine done %d \n", engine.workerCount)
			if engine.workerCount <= 0 {
				// Done
				if len(engine.foundMatch) == 0 {
					fmt.Printf("\n[FileEngine] scheduler() INFO: end scheduler() ,exitstatus=%s \n", "workerDone")
					return
				}
				fmt.Printf("[FileEngine] scheduler() DEBUG: wait channel [foundMatch] clear, left %d \n", len(engine.foundMatch))
				for {
					subRes := <-engine.foundMatch
					*r = append(*r, subRes)
					if len(engine.foundMatch) == 0 {
						return
					}
				}
			}

		// hit select
		case subRes := <-engine.foundMatch:
			//fmt.Printf("[FileEngine] scheduler() DEBUG: found match ,path = %s,sum = %d \n", subRes,len(*r))
			*r = append(*r, subRes)
		}

	}
}

func (engine *FileEngine) search(path string, isRoot bool) {
	files, _ := ioutil.ReadDir(path)
	//fmt.Printf("path : %s, count %d \n",path,len(files))
	for _, file := range files {
		disc := fmt.Sprintf("[light_blue]][%d/%d][reset] Searching... (%d hit)", engine.workerCount, engine.maxWorker, len(*r))
		bar.Describe(disc)
		_ = bar.Add(1)

		// in Dir
		if file != nil && file.IsDir() && !isInArray(&engine.ignorePath, file.Name()) {
			nextDir := path + file.Name() + "/"
			//fmt.Printf("nextDir : %s \n",nextDir)
			engine.m.Lock()
			newGoroutine := engine.maxWorker > engine.workerCount
			engine.m.Unlock()
			if newGoroutine {
				engine.searchRequest <- nextDir
			} else {
				engine.search(nextDir, false)
			}

			// Hit
		} else if isInSuffix(&engine.targetType, file.Name()) {
			if len(engine.foundMatch) > 15 {
				time.Sleep(time.Microsecond * 10)
			}
			engine.foundMatch <- path + file.Name()
		}
	}
	defer func() {
		if isRoot {
			engine.workerDone <- true
		}
	}()
}

// ------------------ Simple Iterator ------------------

// NewEngine rt
func NewEngine(targetType []string, ignorePath []string, result *[]string) *FileEngine {
	engine := FileEngine{}
	if !isAllEmpty(&targetType) && !isAllEmpty(&ignorePath) {
		engine.targetType = targetType
		engine.ignorePath = ignorePath
	} else {
		engine.targetType = []string{".jpg", ".png"}
		engine.ignorePath = []string{".git"}
		fmt.Printf("[FileEngine] CustomEngine() WARN: `targetType` and `ignorePath` is empty, use default")
	}

	engine.workerCount = 0
	engine.maxWorker = 8
	engine.searchRequest = make(chan string)
	engine.workerDone = make(chan bool)
	engine.foundMatch = make(chan string)

	r = result

	return &engine
}

// StartSimpleIterator :use filepath.Walk to traverse path
func (engine *FileEngine) StartSimpleIterator(path string) (*[]string, int) {
	startTime := time.Now()
	fmt.Printf("[FileEngine] %s StartSimpleIterator() INFO: Starting index on %s \n", startTime.Format("20060102 15:04:05"), path)
	engine.workerCount = 1
	// 修复 path 中的分界符
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	engine.walkerSearch(path)

	fmt.Printf("\n[FileEngine] StartSimpleIterator() End index ,TimeConsum %s, found %d \n", time.Since(startTime).String(), len(*r))
	return r, len(*r)
}

// 顺序遍历
func (engine *FileEngine) walkerSearch(path string) {
	bar := basicBar(-1, "[cyan][][reset] Searching...")
	defer func() {
		_ = bar.Finish()
	}()
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		disc := fmt.Sprintf("[cyan][][reset] Searching... (%d hit)", len(*r))
		bar.Describe(disc)
		if f != nil && !f.IsDir() {
			if isInSuffix(&engine.targetType, f.Name()) {
				*r = append(*r, f.Name())
			}
		}
		return nil
	})
	if err != nil {
		fmt.Print(err.Error())
	}
}

// ------------------ Library ------------------

func basicBar(length int, disc string) *progressbar.ProgressBar {
	return progressbar.NewOptions(length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionThrottle(500*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetDescription(disc),
	)
}

// 判断目标字符串是否是在数组中
func isInArray(list *[]string, s string) (isIn bool) {

	if len(*list) == 0 {
		return false
	}
	isIn = false
	for _, f := range *list {
		if f == s {
			isIn = true
			break
		}
	}
	return isIn
}

// 判断后缀名是否命中
func isInSuffix(list *[]string, s string) (isIn bool) {
	ext := ext(s)
	isIn = false
	for _, f := range *list {
		if ext == f {
			isIn = true
			break
		}
	}

	return isIn
}

// Ext 返回文件后缀名
func ext(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

// 判断数组各元素是否是空字符串或空格
func isAllEmpty(list *[]string) (isEmpty bool) {
	if len(*list) == 0 {
		return true
	}
	isEmpty = true
	for _, f := range *list {
		if strings.TrimSpace(f) != "" {
			isEmpty = false
			break
		}
	}
	return isEmpty
}
