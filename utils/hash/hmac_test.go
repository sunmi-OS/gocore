package hash

import (
	"fmt"
	"testing"
)

func TestHmacSHA256(t *testing.T) {

	fmt.Println(HmacSHA256("41414141", "421414124141241"))

}
