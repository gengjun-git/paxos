package acceptor

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_syncFlag(t *testing.T) {
	file, err := os.OpenFile("/home/jun/Desktop/a", os.O_SYNC|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Failed()
	}

	tStart := time.Now()
	for i := 0; i < 100; i++ {
		if err := file.Truncate(0); err != nil {
			t.Failed()
		}
		if _, err := file.Seek(0, 0); err != nil {
			t.Failed()
		}
		if _, err := file.WriteString(fmt.Sprintf("%d", i)); err != nil {
			t.Failed()
		}
	}
	tUsed := time.Now().Sub(tStart)
	fmt.Printf("total [%v], evarage[%v]", tUsed, tUsed/100)
}

func Test_syncFunc(t *testing.T) {
	file, err := os.OpenFile("/home/jun/Desktop/a", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Failed()
	}

	tStart := time.Now()
	for i := 0; i < 100; i++ {
		if err := file.Truncate(0); err != nil {
			t.Failed()
		}
		if _, err := file.Seek(0, 0); err != nil {
			t.Failed()
		}
		if _, err := file.WriteString(fmt.Sprintf("%d", i)); err != nil {
			t.Failed()
		}
		if err := file.Sync(); err != nil {
			t.Failed()
		}
	}
	tUsed := time.Now().Sub(tStart)
	fmt.Printf("total [%v], evarage[%v]", tUsed, tUsed/100)
}
