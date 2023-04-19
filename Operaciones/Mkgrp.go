package operaciones

import (
	"Proyecto2/Estructuras"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Mkgrp(parameters Estructuras.ParamStruct) {
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()
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
	var Grupos []string

	for _, i := range Splited {
		if strings.Contains(i, ",G,") {
			if strings.Contains(i, ",G,"+parameters.Nombre) {
				fmt.Println("Ya existe el grupo")
				return
			}
			Grupos = append(Grupos, i)
		}
	}

	Last := Grupos[len(Grupos)-1]
	n := Last[0]
	num, err := strconv.ParseInt(string(n), 10, 64)
	NewGrp := string(fmt.Sprint(num+1)) + ",G," + parameters.Nombre + "\n"
	WrtArchivoMkfs(SuperBlock, "users.txt", dsk, NewGrp, false)
}

func Rmgrp(parameters Estructuras.ParamStruct) {
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()
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
		if strings.Contains(Splited[i], ",G,"+parameters.Nombre) {
			Splited[0] = "0"
			break
		}
	}
	Novo := ""
	for _, i := range Splited {
		Novo += i
	}
	WrtArchivoMkfs(SuperBlock, "users.txt", dsk, Novo, true)
}
