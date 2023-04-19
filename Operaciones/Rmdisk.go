package operaciones

import (
	"Proyecto2/Estructuras"
	"bufio"
	"fmt"
	"os"
)

func Rmdisk(parameters Estructuras.ParamStruct) {
	fullpath := parameters.Direccion
	if !existeDisco(fullpath) {
		fmt.Println("---NO SE PUEDE ELIMINAR EL DISCO DEBIDO A QUE NO EXISTE---")
		return
	}

	for {
		fmt.Println("Esta seguro que desea borrar el disco " + fullpath + "?")
		fmt.Println("1) Si")
		fmt.Println("2) No")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		entrada := scanner.Text()

		if entrada == "1" {
			fmt.Println("Se confirmo que se desea eliminar el disco " + fullpath + ", se procedera a elminarlo")
			err := os.Remove(fullpath) // remove a single file
			if err != nil {
				fmt.Println(err)
			}
			return
		} else if entrada == "2" {
			fmt.Println("Se cancelara la eliminacion")
			return
		} else {
			fmt.Println("Ingrese una opcion valida")
		}
	}

}
