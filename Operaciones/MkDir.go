package operaciones

import (
	"Proyecto2/Estructuras"
	"fmt"
	"os"
	"strings"
)

func Mkdir(parameters Estructuras.ParamStruct) {
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()
	StartPoint := Uss.Startpoint
	path := Uss.Path
	VirtualPath := parameters.Direccion

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	SuperBlock := ReadSBlock(Estructuras.Sblock{}, StartPoint, dsk)
	pathSeparado := strings.Split(VirtualPath, "/")
	if pathSeparado[0] == "" {
		pathSeparado = append(pathSeparado[1:])
	}

	SrchInodo(0, dsk, pathSeparado, SuperBlock, true, 0)
}
