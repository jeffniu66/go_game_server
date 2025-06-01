package tableconfig

import (
	"fmt"
	"testing"
)

func TestReadTable(t *testing.T) {
	fmt.Println(TaskConfigs)
	ReadTable()
	fmt.Println(TaskConfigs.TaskList)
	fmt.Println(TaskConfigs.TaskMap)
}
