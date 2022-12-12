package main

import (
	"fmt"

	"github.com/velotio-ajaykumbhar/controller/opration"
)

func init() {
	opration.New()
}

func main() {
	for {
		var input int
		fmt.Printf("Choose opration\n 1) Deployment\n 2) Service\n 3) ConfigMap \n Input : ")
		fmt.Scanf("%d", &input)

		switch input {
		case 1:
			opration.Deployment()
		case 2:
			opration.Service()
		case 3:
			opration.ConfigMap()
		default:
			fmt.Println("Invalid choice")
		}
	}
}
