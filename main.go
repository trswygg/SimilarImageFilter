package main

import (
	imagefilter "SimilarImageFilter/ImageFilter"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"runtime"
)

type setting struct {
	WorkDir    string
	byHash     string // a,d,p
	recursive  bool

	Width    int
	Height   int
	Dist     int
	Thread   int
	MaxImage int
}
var set = setting{}

func init() {
	var defaultWorkDir string
	if runtime.GOOS == "linux" {
		defaultWorkDir = "/mnt/i/"
	} else if runtime.GOOS == "windows" {
		defaultWorkDir = "I:"
	}

	flag.StringVar(&set.WorkDir, "d", defaultWorkDir, "Work dir for ImageSearch engine")
	flag.StringVar(&set.byHash, "hash", "adp", "Which hash func to use(AverageHash -a,DifferenceHash -d,PerceptionHash,-p)")
	flag.BoolVar(&set.recursive, "r", true, " Recursive search in work dir")

	flag.IntVar(&set.Width,"w",16,"Width of resized image")
	flag.IntVar(&set.Height,"h",16,"Height of resized image")
	flag.IntVar(&set.Dist,"dist",50,"Hamming Distance")
	flag.IntVar(&set.Thread,"thread",4,"Max work thread")
	flag.IntVar(&set.MaxImage,"MaxImage",100,"Max image (to shorten the comparison time)")
}

func main() {
	flag.Parse()

	comp()
}

func comp() {
	filter := imagefilter.NewImageFilter(
		imagefilter.SetWidth(set.Width),
		imagefilter.SetHeight(set.Height),
		imagefilter.SetDist(set.Dist),
		imagefilter.SetSource(imagefilter.UseFileengine),
		imagefilter.SetHash(set.byHash),
		imagefilter.SetThread(set.Thread),
		imagefilter.SetMax(set.MaxImage),
	)
	filter.Index(set.WorkDir)
	imagefilter.ShowMemoryInfo()
	filter.Contras()
	filter.Show(true,func (g imagefilter.SimilarPicture) {
		fmt.Printf("%s %v By:%c \n", g.A, g.B, g.By)
	})
}
