package arfc

import "os"
import "fmt"
import "github.com/sashakoshka/arf"

func main () {
	if len(os.Args) != 2 {
		fmt.Println("XXX please specify module path and output path")
	}

	inPath  := os.Args[0]
	outPath := os.Args[1]
	file := os.OpenFile(outPath, os.WRONLY | os.CREATE, 0655)

	err := arf.CompileModule(inPath, outFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
