package design

import (
	"fmt"

	"github.com/Brum3ns/firefly/pkg/firefly/info"
)

func Banner() {
	fmt.Printf(`                             
   / __7 o          / _7/¯7  
  / _7 /¯7 /¯_7¯-_)/ _7/ /\¯\/7
 /_/  /_/ /_/ \__7/_/ /_/  ) /
                          /_/
                                (%s)
   By: @yeswehack : Brumens
`, (COLOR.GREY + info.VERSION + COLOR.WHITE))
}

func Disclaimer() {
	fmt.Println(ICON.AWARE + " Stay ethical. The creator of the tool is not responsible for any misuse or damage.")
}
