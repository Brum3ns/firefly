package design

import "fmt"

func Loadingbar_Verify(nr, total int) {
	percent := (int)((nr+1)*100) / total

	fmt.Printf("\r%s Verify target behaviour: %d%%", Info, percent)
	if percent >= 100 {
		fmt.Println()
	}
}
