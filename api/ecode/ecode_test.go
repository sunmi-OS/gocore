package ecode

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestName(t *testing.T) {

	err := New(255, gorm.ErrRecordNotFound)
	//err := New(255, errors.New("record not found"))
	fmt.Println(err)

	fmt.Println(Transform(gorm.ErrRecordNotFound))

}
