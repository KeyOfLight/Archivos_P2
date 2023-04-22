package operaciones

import (
	"Proyecto2/Estructuras"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

func Lgn(parameters Estructuras.ParamStruct) {
	var Listado = Estructuras.ListaMontados{}
	Listado = (&Listado).GetLista()
	PMontada := false
	var StartPoint int64
	path := ""

	for _, i := range Listado.Montado {
		if parameters.Nombre == i.Id {
			path = i.Path
			StartPoint = i.StartPoint
			PMontada = true
			break
		}
	}

	if !PMontada {
		fmt.Println("La particion deseada no esta montada")
	}

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	SuperBlock := Estructuras.Sblock{}
	SuperBlock = ReadSBlock(SuperBlock, StartPoint, dsk)

	bl_archivos := Estructuras.BloqueArchivos{}
	bl_archivos = ReadBloqueArchivo(bl_archivos, SuperBlock.S_block_start+int64(unsafe.Sizeof(Estructuras.BloqueArchivos{})), dsk)

	Bloquear := strings.Split(string(bl_archivos.B_content[:]), "\n")
	UsuarioInf := ReturnUsser(Bloquear)

	for _, i := range UsuarioInf {
		var Name = [11]byte{}
		copy(Name[:], []byte(parameters.User))
		var Pass = [11]byte{}
		copy(Pass[:], []byte(parameters.Pwd))
		if i.Uss == Name {
			if i.Pass == Pass {
				var Uss = Estructuras.User{}
				i.Path = path
				i.Startpoint = StartPoint
				(&Uss).Loguear(i)
				return
			} else {
				fmt.Println("No se pudo ingresar sesion")
			}
		}
	}
}

func LgOut() {
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()

	if Uss.Uss == [11]byte{} {
		fmt.Println("No hay ninguna sesion iniciada")
		return
	} else {
		(&Uss).Logout()
		return
	}
}

func ReturnUsser(Datos []string) []Estructuras.User {

	Usuario := Estructuras.User{}
	List := []Estructuras.User{}

	for _, i := range Datos {
		if strings.Contains(i, ",U,") {
			Splited := strings.Split(i, ",")
			for _, h := range Datos {
				if strings.Contains(h, ",G,"+Splited[3]) {
					Usuario.GUID = h[0]
					break
				}
			}
			Uid := Splited[0]
			Usuario.UID = Uid[0]
			copy(Usuario.Grupo[:], []byte(Splited[2]))
			copy(Usuario.Uss[:], []byte(Splited[3]))
			copy(Usuario.Pass[:], []byte(Splited[4]))
			List = append(List, Usuario)
		}
	}

	return List
}
