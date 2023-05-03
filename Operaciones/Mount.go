package operaciones

import (
	"Proyecto2/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
)

func MountPart(parameters Estructuras.ParamStruct) {
	fullpath := parameters.Direccion
	if !existeDisco(fullpath) {
		fmt.Println("---NO SE PUEDE ELIMINAR EL DISCO DEBIDO A QUE NO EXISTE---")
		return
	}

	mbr := OpenMBR(fullpath)
	PartExt := Estructuras.Particion{}

	var Listado = Estructuras.ListaMontados{}
	Listado = (&Listado).GetLista()

	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_type == 'e' {
			PartExt = mbr.Mbr_partition[i]

		}
		if string(mbr.Mbr_partition[i].Part_name[:len(parameters.Nombre)]) == (parameters.Nombre) {

			(&Listado).Montar(parameters.Nombre, fullpath, mbr.Mbr_partition[i].Part_start, mbr.Mbr_partition[i].Part_size, (mbr.Mbr_partition[i].Part_type))
			fmt.Println("Se ha montado la particion")
			fmt.Println("Id: " + parameters.Nombre)
			return
		}
	}

	if PartExt.Part_status == '1' {
		MountExt(PartExt.Part_start, parameters, fullpath)
		return
	}

}

func MountExt(startpoint int64, parameters Estructuras.ParamStruct, fullpath string) {
	tempEbr := Estructuras.EBR{}

	dsk, err := os.OpenFile(fullpath, os.O_RDWR, 0777)
	defer dsk.Close()
	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	var sizeMbr int64 = int64(unsafe.Sizeof(tempEbr))
	dsk.Seek(0, 0)
	dsk.Seek(startpoint, 0)

	var Listado = Estructuras.ListaMontados{}
	(&Listado).GetLista()

	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &tempEbr)
	if err2 != nil {
		log.Fatal("binary.Read failed", err2)
	}

	if string(tempEbr.Part_name[:len(parameters.Nombre)]) == (parameters.Nombre) {
		(&Listado).Montar(parameters.Nombre, parameters.Direccion, startpoint, tempEbr.Part_size, 'E')
		fmt.Println("Se ha montado la particion")
		return
	}

	if tempEbr.Part_next != -1 {
		MountExt(startpoint, parameters, fullpath)
	} else {
		fmt.Println("No se encontro La particion deseada.")
		return
	}
}
