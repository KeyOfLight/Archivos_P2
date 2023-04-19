package operaciones

import (
	"Proyecto2/Estructuras"
	"fmt"
	"os"
	"strings"
)

func MkUsr(parameters Estructuras.ParamStruct) {
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()
	Root := [11]byte{}
	copy(Root[:], "root")
	exstGrp := false

	if Uss.Uss != Root {
		fmt.Println("Solo el usuario root puede crear usuarios")
		return
	}
	StartPoint := Uss.Startpoint
	path := Uss.Path

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	SuperBlock := ReadSBlock(Estructuras.Sblock{}, StartPoint, dsk)
	Info := LeerArchivoMkfs(SuperBlock, "users.txt", dsk)
	Splited := strings.Split(Info, "\n")
	var n byte
	for _, i := range Splited {
		if strings.Contains(i, ",G,") {
			if strings.Contains(i, ",G,"+parameters.Direccion) {
				exstGrp = true
				Last := i
				n = Last[0]
			}

		}
		if strings.Contains(i, ",U,") {
			if strings.Contains(i, ",U,"+parameters.Nombre) {
				fmt.Println("Ya existe el usuario")
				return
			}
		}
	}

	if !exstGrp {
		fmt.Println("No existe el grupo")
		return
	}

	NewGrp := string(n) + ",U," + parameters.Direccion + "," + parameters.Nombre + "," + parameters.Pwd + "\n"
	WrtArchivoMkfs(SuperBlock, "users.txt", dsk, NewGrp, false)
}

func RmUsr(parameters Estructuras.ParamStruct) {
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()
	Root := [11]byte{}
	copy(Root[:], "root")

	if Uss.Uss != Root {
		fmt.Println("Solo el usuario root puede crear usuarios")
		return
	}
	StartPoint := Uss.Startpoint
	path := Uss.Path

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	SuperBlock := ReadSBlock(Estructuras.Sblock{}, StartPoint, dsk)
	Info := LeerArchivoMkfs(SuperBlock, "users.txt", dsk)
	Splited := strings.Split(Info, "\n")

	for i := 0; i < len(Splited); i++ {
		if strings.Contains(Splited[i], ",U,"+parameters.Nombre) {
			Splited[0] = "0"
			return
		}
	}
	Novo := ""
	for _, i := range Splited {
		Novo += i
	}
	WrtArchivoMkfs(SuperBlock, "users.txt", dsk, Novo, true)
}
