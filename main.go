package main

import (
	"Proyecto2/Analizador"
	"bufio"
	"fmt"
	"os"
)

func main() {

	intContador := 0

	for {
		fmt.Println("")
		fmt.Println("-  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -")
		fmt.Println("-----------------------[MIA] Proyecto 2-----------------------")
		fmt.Println("------------------Diego Andre Gomez 201908327------------------")
		fmt.Println("")
		fmt.Print("----Ingrese comando: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		entrada := scanner.Text()

		Analizador.Analizar(entrada)

		fmt.Println("\n-------------------------------------------")
		fmt.Println("---------------------------------------------")
		fmt.Println(" ")

		intContador++
	}

}
