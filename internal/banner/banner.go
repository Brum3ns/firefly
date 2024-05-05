package banner

import (
	"fmt"

	"github.com/Brum3ns/firefly/internal/version"
	"github.com/Brum3ns/firefly/pkg/design"
)

func Banner() {
	fmt.Printf(`                             
   / __7 o          / _7/¯7  
  / _7 /¯7 /¯_7¯-_)/ _7/ /\¯\/7
 /_/  /_/ /_/ \__7/_/ /_/  ) /
                          /_/
                                (%s)
   By: @yeswehack : Brumens

`, (design.COLOR.GREY + version.VERSION + design.COLOR.WHITE))
}

func Disclaimer() {
	fmt.Println(design.ICON.AWARE + " Stay ethical. The creator of the tool is not responsible for any misuse or damage.")
}
