package Analizador

import (
	"Proyecto2/Estructuras"
	operaciones "Proyecto2/Operaciones"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func Analizar(Cadena string) {

	var entradacmd string
	entradacmd = Cadena
	entradacmd = strings.ToLower(entradacmd)
	Separado := false

	VectorEntrada := strings.Split(entradacmd, " ")
	var Aux []string
	path := ""
	for i := 0; i < len(VectorEntrada); i++ {
		param := VectorEntrada[i]
		if strings.Contains(param, ">path=\"") {
			path += strings.Replace(param, "\"", "", -1)
			Separado = true
		} else if Separado {
			path += " " + strings.Replace(param, "\"", "", -1)
			if strings.Contains(param, "\"") {
				Separado = false
				Aux = append(Aux, path)
			}
		} else {
			Aux = append(Aux, param)
		}

	}
	var comadno string
	var parametros []string
	for i := 0; i < len(Aux); i++ {
		if i == 0 {
			comadno = Aux[i]
		} else {
			parametros = append(parametros, Aux[i])
		}
	}

	ReconocerComando(comadno, parametros)

}

func ReconocerComando(comando string, Parametros []string) {

	cmd := Estructuras.ParamStruct{}
	if comando == "execute" {
		cmd.Nombre = "execute"

		if len(Parametros) == 0 {
			fmt.Println("Debe ingresar un path para poder ejecutar este comando")
			return
		}

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				cmd.Direccion = param
			}
		}

		Executar(cmd.Direccion)

	} else if comando == "mkdisk" {
		cmd.Nombre = "mkdisk"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">size=") {
				param = strings.ReplaceAll(param, ">size=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Tam = param
			} else if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Direccion = param
			} else if strings.Contains(param, ">fit=") {
				param = strings.ReplaceAll(param, ">fit=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Fit = param
			} else if strings.Contains(param, ">unit=") {
				param = strings.ReplaceAll(param, ">unit=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Unit = param
			}
		}
		operaciones.MakeDisk(cmd)
	} else if comando == "rmdisk" {
		cmd.Nombre = "rmdisk"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Direccion = param
			}
		}
		operaciones.Rmdisk(cmd)
	} else if comando == "fdisk" {
		cmd.Nombre = "fdsik"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">size=") {
				param = strings.ReplaceAll(param, ">size=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Tam = param
			} else if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Direccion = param
			} else if strings.Contains(param, ">name=") {
				param = strings.ReplaceAll(param, ">name=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			} else if strings.Contains(param, ">unit=") {
				param = strings.ReplaceAll(param, ">unit=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Unit = param
			} else if strings.Contains(param, ">type=") {
				param = strings.ReplaceAll(param, ">type=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Tipo = param
			} else if strings.Contains(param, ">fit=") {
				param = strings.ReplaceAll(param, ">fit=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Fit = param
			} else if strings.Contains(param, ">delete=") {
				param = strings.ReplaceAll(param, ">delete=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Delete = param
			} else if strings.Contains(param, ">add=") {
				param = strings.ReplaceAll(param, ">add=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Add = param
			}
		}
		operaciones.CrearParticion(cmd)
	} else if comando == "mount" {
		cmd.Nombre = "mount"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Direccion = param
			} else if strings.Contains(param, ">name=") {
				param = strings.ReplaceAll(param, ">name=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			}
		}
		operaciones.MountPart(cmd)
	} else if comando == "mkfs" {
		cmd.Nombre = "mkfs"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">id=") {
				param = strings.ReplaceAll(param, ">id=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			}
		}
		operaciones.MakeFs(cmd)
	} else if comando == "login" {
		cmd.Nombre = "login"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">id=") {
				param = strings.ReplaceAll(param, ">id=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			} else if strings.Contains(param, ">usuario=") {
				param = strings.ReplaceAll(param, ">usuario=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.User = param
			} else if strings.Contains(param, ">password=") {
				param = strings.ReplaceAll(param, ">password=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Pwd = param
			}
		}
		operaciones.Lgn(cmd)
	} else if comando == "logout" {
		cmd.Nombre = "logout"
		operaciones.Lgn(cmd)

	} else if comando == "mkgrp" {
		cmd.Nombre = "mkgrp"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">name=") {
				param = strings.ReplaceAll(param, ">name=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			}
		}
		operaciones.Mkgrp(cmd)
	} else if comando == "rmgrp" {
		cmd.Nombre = "rmgrp"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">name=") {
				param = strings.ReplaceAll(param, ">name=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			}
		}
		operaciones.Rmgrp(cmd)
	} else if comando == "mkfile" {
		cmd.Nombre = "mkfile"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Direccion = param
			} else if strings.Contains(param, ">r") {
				param = strings.ReplaceAll(param, ">r", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Pwd = "s"
			} else if strings.Contains(param, ">name=") {
				param = strings.ReplaceAll(param, ">name=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			} else if strings.Contains(param, ">size=") {
				param = strings.ReplaceAll(param, ">size=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Size = param
			} else if strings.Contains(param, ">cont=") {
				param = strings.ReplaceAll(param, ">cont=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Tam = param
			}
		}
		operaciones.Mkfile(cmd)
	} else if comando == "mkdir" {
		cmd.Nombre = "mkdir"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Direccion = param
			} else if strings.Contains(param, ">r") {
				param = strings.ReplaceAll(param, ">r", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Pwd = "s"
			}
		}
		operaciones.Mkdir(cmd)
	} else if strings.Contains(comando, "pause") {

		fmt.Println("--------------PAUSA-----------------")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

	} else if strings.Contains(comando, "#") {
		fmt.Print(comando)
		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			fmt.Print(" " + param)
		}
		fmt.Println(" ")

	} else if comando == "rep" {
		cmd.Nombre = "rep"

		for i := 0; i < len(Parametros); i++ {
			param := Parametros[i]
			if strings.Contains(param, ">path=") {
				param = strings.ReplaceAll(param, ">path=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Direccion = param
			} else if strings.Contains(param, ">name=") {
				param = strings.ReplaceAll(param, ">name=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Nombre = param
			} else if strings.Contains(param, ">id=") {
				param = strings.ReplaceAll(param, ">id=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Id = param
			} else if strings.Contains(param, ">ruta=") {
				param = strings.ReplaceAll(param, ">ruta=", "")
				param = strings.ReplaceAll(param, "\"", "")
				cmd.Pwd = param
			}
		}
		operaciones.Reportes(cmd)
	}
}

func Executar(Dir string) {
	Datos, err := ioutil.ReadFile(Dir)

	if err != nil {
		log.Fatal(err)
	}
	datosComoString := string(Datos)
	SeparadoLineas := strings.Split(datosComoString, "\n")
	for _, linea := range SeparadoLineas {
		Analizar(linea)
	}

}
