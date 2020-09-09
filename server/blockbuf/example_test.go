package blockbuf

import (
		"os"
		"fmt"
		"testing"
	   )

func TestExample(t *testing.T) {
	bf := New(os.Stdin, os.Stdout)

	fmt.Println(": ", bf.String())
}
