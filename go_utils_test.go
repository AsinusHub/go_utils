package go_utils

import (
	"fmt"
)

func ExampleCreateJson(){
	tmp:=ErrorResponse{
		Status:"error",
		Description:"This example cannot work without created struct",
	}

	jsn := CreateJson(tmp)
	fmt.Print(string(jsn))
	// Output:{"status":"error","description":"This example cannot work without created struct"}
}
