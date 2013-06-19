package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	cssDir := "\\theme\\images\\sdf\\"
	imgDir := "\\theme\\images\\sdf\\"
	rel, _ := filepath.Rel(cssDir , imgDir )
	rel = strings.Replace(rel,"/\\","/",-1)
	rel = strings.Replace(rel,"\\","/",-1)
	fmt.Println(rel)
}