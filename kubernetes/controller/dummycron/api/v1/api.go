package v1

import (
	"fmt"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func testTime() {
	now := metaV1.Now()
	fmt.Print(now)
}
