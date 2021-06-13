package imagefilter

import (
	"SimilarImageFilter/ImageFilter/Library/go-libjpeg/jpeg"
	"bytes"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// UTIL

func splitArray(arr []string, num int) [][]string {
	max := len(arr)
	if max < num {
		return nil
	}
	var segments = make([][]string, 0)
	quantity := max / num
	end := 0
	for i := 1; i <= num; i++ {
		qu := i * quantity
		if i != num {
			segments = append(segments, arr[i-1+end:qu])
		} else {
			segments = append(segments, arr[i-1+end:])
		}
		end = qu - i
	}
	return segments
}

// exists returns whether the given file or directory exists or not
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

type MemInfo struct {
	Total       uint64
	Used        uint64
	UsedPercent float64
}

func GetMemInfo() MemInfo {
	m, _ := mem.VirtualMemory()
	return MemInfo{
		Total:       m.Total,
		Used:        m.Used,
		UsedPercent: m.UsedPercent,
	}
}
func ShowMemoryInfo() {
	m := GetMemInfo()
	sub := [...]string{" ","▍", "▍", "▌", "▌", "▋", "▋", "▊", "▊", "▉"}
	width := 25
	Lcount := int(m.UsedPercent / 100.0 * float64(width))
	mid := sub[int(m.UsedPercent)%10/(100/width)+1]
	Rcount := width - Lcount - 1
	s := fmt.Sprintf("[%s%s%s]", strings.Repeat("█", Lcount), mid, strings.Repeat(" ", Rcount))
	fmt.Printf("[INFO]total memory used [%5s / %5s]MB %s %.1f%% \n", strconv.FormatUint(m.Used>>20, 10), strconv.FormatUint(m.Total>>20, 10), s, m.UsedPercent)
}

func getCpuUse() float64 {
	percent, _ := cpu.Percent(time.Millisecond*500, false)
	return percent[0]
}

func basicBar(length int, disc string) *progressbar.ProgressBar {
	return progressbar.NewOptions(length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionThrottle(500*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(51),
		progressbar.OptionSetDescription(disc),
	)
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

//func encode(path string, file *os.File, rgba *image.RGBA) error {
//	if strings.HasSuffix(path, "jpg") || strings.HasSuffix(path, "jpeg") {
//		jpeg.Encode(file, rgba, nil)
//	} else if strings.HasSuffix(path, "png") {
//		png.Encode(file, rgba)
//	} else {
//		return fmt.Errorf("不支持的图片格式")
//	}
//}

func decodeImg(path string) (image.Image, error) {
	var i image.Image
	var err error
	op := jpeg.DecoderOptions{
		ScaleTarget:            image.Rectangle{},
		DCTMethod:              jpeg.DCTFloat,
		DisableFancyUpsampling: false,
		DisableBlockSmoothing:  false,
	}
	f, _ := ioutil.ReadFile(path)
	buf := bytes.NewBuffer(f)
	if ext(path) == ".jpg" {
		// file -> image.YCbCr，decoding partial images
		i, err = jpeg.Decode(buf, &op)
		if err != nil {
			i, _, err = image.Decode(buf)
		}
	} else {
		i, _, err = image.Decode(buf)
	}
	return i, err
}
