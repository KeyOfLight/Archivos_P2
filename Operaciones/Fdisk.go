package operaciones

import (
	"Proyecto2/Estructuras"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func CrearParticion(parameters Estructuras.ParamStruct) {

	if parameters.Delete == "full" {
		EliminarParticion(parameters)
		return
	}

	fullpath := parameters.Direccion

	size, _ := strconv.ParseInt(parameters.Tam, 10, 64)
	size = RealSize(size, parameters.Unit)

	if !existeDisco(fullpath) {
		fmt.Println("---No se encontro el disco buscado---")
		return
	}

	mbr := Estructuras.MBR{}

	var sizeMbr int64 = int64(unsafe.Sizeof(mbr))

	dsk, err := os.OpenFile(fullpath, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		return
	}

	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &mbr)
	if err2 != nil {
		log.Fatal("binary.Read failed", err)
	}

	CurrentFit := VerFit(parameters.Fit)
	startpoint := 0

	if parameters.Tipo == "p" || parameters.Tipo == "" || parameters.Tipo == "e" {
		if CurrentFit == 3 {
			startpoint = GetPositionFF(mbr, size)
			parameters.Fit = "b"
		} else if CurrentFit == 2 {
			startpoint = GetPositionFF(mbr, size)
			parameters.Fit = "w"
		} else if CurrentFit == 1 {
			startpoint = GetPositionFF(mbr, size)
			parameters.Fit = "f"
		}
	}

	if parameters.Tipo == "p" || parameters.Tipo == "" {
		parameters.Tipo = "p"
		filledsize := tamparticiones(mbr)
		if size < (mbr.Mbr_tamano - int64(filledsize)) {
			EmptPart := Estructuras.MBR{}
			for i := 0; i < 4; i++ {
				if parameters.Nombre == string(mbr.Mbr_partition[i].Part_name[:len(parameters.Nombre)]) {
					fmt.Println("No se puede crear una particion con un nombre ya utilizado")
					return
				}
			}

			for i := 0; i < 4; i++ {
				if mbr.Mbr_partition[i].Part_name == EmptPart.Mbr_partition[0].Part_name {
					mbr.Mbr_partition[i] = ActivarParticion(mbr.Mbr_partition[i], parameters, size, int64(startpoint))
					dsk.Seek(0, 0)
					var bufferControl bytes.Buffer
					binary.Write(&bufferControl, binary.BigEndian, &mbr)
					_, err = dsk.Write(bufferControl.Bytes())
					PrintPartAct(mbr.Mbr_partition[i])
					return

				}
			}
		} else {
			fmt.Println("La particion debe ser menor al tam total del disco")
			return
		}
	} else if parameters.Tipo == "e" {
		filledsize := tamparticiones(mbr)
		if size < (mbr.Mbr_tamano - int64(filledsize)) {
			EmptPart := Estructuras.MBR{}
			for i := 0; i < 4; i++ {
				if parameters.Nombre == string(mbr.Mbr_partition[i].Part_name[:len(parameters.Nombre)]) {
					fmt.Println("No se puede crear una particion con un nombre ya utilizado")
					return
				}

				if mbr.Mbr_partition[i].Part_type == 'e' {
					fmt.Println("Solo se puede crear una particion de tipo extendida")
					return
				}
			}

			for i := 0; i < 4; i++ {
				if mbr.Mbr_partition[i].Part_name == EmptPart.Mbr_partition[0].Part_name {
					mbr.Mbr_partition[i] = ActivarParticion(mbr.Mbr_partition[i], parameters, size, int64(startpoint))
					dsk.Seek(0, 0)
					var bufferControl bytes.Buffer
					binary.Write(&bufferControl, binary.BigEndian, &mbr)
					_, err = dsk.Write(bufferControl.Bytes())
					PrintPartAct(mbr.Mbr_partition[i])
					return

				}
			}
		} else {
			fmt.Println("La particion debe ser menor al tam total del disco")
			return
		}
	} else if parameters.Tipo == "l" {
		filledsize := tamparticiones(mbr)
		if size < (mbr.Mbr_tamano - int64(filledsize)) {
			EmptPart := Estructuras.EBR{}
			PartExt := Estructuras.Particion{}
			for i := 0; i < 4; i++ {
				if parameters.Nombre == string(mbr.Mbr_partition[i].Part_name[:len(parameters.Nombre)]) {
					fmt.Println("No se puede crear una particion con un nombre ya utilizado")
					return
				}

				if mbr.Mbr_partition[i].Part_type == 'e' {
					PartExt = mbr.Mbr_partition[i]
				}
			}

			if PartExt.Part_status != '1' {
				fmt.Println("Para crear una particion logica es necesario que exista una particion extendida")
				return
			}
			if parameters.Nombre == "logicp1" {
				fmt.Println("ya")
			}
			if SearchNameEBR(PartExt, parameters.Nombre, parameters) {
				fmt.Println("No se puede crear una particion logica con un nombre ya utilizado")
				return
			}

			libre := (PartExt.Part_size - CalcularSizeExt(PartExt.Part_start, parameters, 0)) - size

			if libre < 0 {
				fmt.Println("El espacio restante en el disco no es suficiente para almacenar esta particion logica")
			}

			if CurrentFit == 3 {
				EmptPart = GetPositionEBRFF(PartExt, size, parameters)
				parameters.Fit = "b"
			} else if CurrentFit == 2 {
				EmptPart = GetPositionEBRFF(PartExt, size, parameters)
				parameters.Fit = "w"
			} else if CurrentFit == 1 {
				EmptPart = GetPositionEBRFF(PartExt, size, parameters)
				parameters.Fit = "f"
			}

			EmptPart = ActivarEBR(parameters, size, EmptPart)

			dsk.Seek((EmptPart.Part_start), 0)
			var bufferControl bytes.Buffer
			binary.Write(&bufferControl, binary.BigEndian, &EmptPart)
			escribirBytes(dsk, bufferControl.Bytes())

			PrintEbrAct(EmptPart)

			return
		} else {
			fmt.Println("La particion debe ser menor al tam total del disco")
			return
		}
	}

}

func CalcularSizeExt(startpoint int64, parameters Estructuras.ParamStruct, ocupado int64) int64 {
	path := parameters.Direccion
	Occupied := ocupado
	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	dsk.Seek((startpoint), 0)
	tempEbr := Estructuras.EBR{}
	var sizeMbr int64 = int64(unsafe.Sizeof(tempEbr))

	dataControl := leerBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	errb := binary.Read(bufferControl, binary.BigEndian, &tempEbr)
	if errb != nil {
		log.Fatal("binary.Read failed", err)
	}

	if tempEbr.Part_status == '1' {
		Occupied += tempEbr.Part_size
		return CalcularSizeExt(tempEbr.Part_next, parameters, Occupied)
	}
	return Occupied
}

func ActivarEBR(parameters Estructuras.ParamStruct, size int64, Ebr Estructuras.EBR) Estructuras.EBR {

	for i := 0; i < len(parameters.Nombre); i++ {
		Ebr.Part_name[i] = parameters.Nombre[i]
	}
	Ebr.Part_fit = parameters.Fit[0]
	Ebr.Part_size = size
	Ebr.Part_status = '1'

	return Ebr

}

func PrintPartAct(Ebr Estructuras.Particion) {

	println("-------------------------------------------------")
	println("Nombre: " + string(Ebr.Part_name[:]))
	println("Part_fit: " + string(Ebr.Part_fit))
	println("Ebr.Part_size: " + fmt.Sprint(Ebr.Part_size))
	println("Ebr.Part_status: " + string(Ebr.Part_status))
	println("-------------------------------------------------")

}

func PrintEbrAct(Ebr Estructuras.EBR) {

	println("-------------------------------------------------")
	println("Nombre: " + string(Ebr.Part_name[:]))
	println("Part_fit: " + string(Ebr.Part_fit))
	println("Ebr.Part_size: " + fmt.Sprint(Ebr.Part_size))
	println("Ebr.Part_status: " + string(Ebr.Part_status))
	println("-------------------------------------------------")

}

func ActivarParticion(part Estructuras.Particion, parameters Estructuras.ParamStruct, size int64, startpoint int64) Estructuras.Particion {

	for i := 0; i < len(parameters.Nombre); i++ {
		part.Part_name[i] = parameters.Nombre[i]
	}
	part.Part_fit = parameters.Fit[0]
	part.Part_size = size
	part.Part_status = '1'
	part.Part_type = parameters.Tipo[0]
	part.Part_start = startpoint

	return part
}

func tamparticiones(mbr Estructuras.MBR) int {

	totalsize := 0
	for _, i := range mbr.Mbr_partition {
		totalsize += int(i.Part_size)
	}

	return totalsize
}

func GetPositionEBRFF(PartExt Estructuras.Particion, size int64, parameters Estructuras.ParamStruct) Estructuras.EBR {

	PosEbr := Estructuras.EBR{}
	path := parameters.Direccion
	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	startpoint := PartExt.Part_start
	defer dsk.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if CalcularSizeExt(startpoint, parameters, 0) == 0 {
		PosEbr.Part_start = startpoint
		PosEbr.Part_next = -1
		return PosEbr
	}

	dsk.Seek((startpoint), 0)
	tempEbr := Estructuras.EBR{}
	var sizeMbr int64 = int64(unsafe.Sizeof(tempEbr))

	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	errb := binary.Read(bufferControl, binary.BigEndian, &tempEbr)
	if errb != nil {
		log.Fatal("binary.Read failed", err)
	}

	var AllArrays []Estructuras.EBR

	pos := 0

	for {
		if tempEbr.Part_status != '1' {
			break
		}
		AllArrays = append(AllArrays, tempEbr)
		pos += 1

		if tempEbr.Part_next == -1 {
			break
		}

		dsk.Seek((tempEbr.Part_next), 0)
		var sizeMbr int64 = int64(unsafe.Sizeof(tempEbr))
		dataControl := ReadBytes(dsk, int(sizeMbr))
		rbufferControl := bytes.NewBuffer(dataControl)
		err2 := binary.Read(rbufferControl, binary.BigEndian, &tempEbr)
		if err2 != nil {
			log.Fatal("binary.Read failed", err)
			fmt.Println(errb)
			break
		}

	}

	n := pos
	var temp []Estructuras.EBR
	postemp := 0

	for i := 1; i < n; i++ {
		for j := n - 1; j >= i; j-- {
			if AllArrays[j-1].Part_start > AllArrays[j].Part_start && AllArrays[j-1].Part_start != 0 && AllArrays[j].Part_start != 0 {
				temp[postemp] = AllArrays[j-1]
				AllArrays[j-1] = AllArrays[j]
				AllArrays[j] = temp[postemp]
				postemp += 1
			}
		}
	}

	if n < 2 {
		nextpos := AllArrays[0].Part_start + AllArrays[0].Part_size
		AllArrays[0].Part_next = nextpos
		Reescribir := AllArrays[0]

		dsk.Seek((Reescribir.Part_start), 0)
		var bufferControl bytes.Buffer
		binary.Write(&bufferControl, binary.BigEndian, &Reescribir)
		escribirBytes(dsk, bufferControl.Bytes())

		PosEbr.Part_start = nextpos
		PosEbr.Part_next = -1
		return PosEbr

	}

	var fin int64
	fin = 0
	for i := 1; i < 4; i++ {
		Actual := AllArrays[i]
		Anterior := AllArrays[i-1]

		if Actual.Part_status == '1' {
			fin = (Anterior.Part_size) + (Anterior.Part_start)
			if fin+(size) < (Actual.Part_start) {
				if Anterior.Part_next == -1 {
					Anterior.Part_next = (fin)
					PosEbr.Part_next = -1
					dsk.Seek((Anterior.Part_start), 0)
					binary.Write(bufferControl, binary.BigEndian, &Anterior)
					escribirBytes(dsk, bufferControl.Bytes())
				} else {
					PosEbr.Part_next = Anterior.Part_next
					Anterior.Part_next = fin
					PosEbr.Part_start = fin
					dsk.Seek((Anterior.Part_start), 0)
					binary.Write(bufferControl, binary.BigEndian, &Anterior)
					escribirBytes(dsk, bufferControl.Bytes())
				}

				return PosEbr

			} else {
				if AllArrays[i].Part_next == -1 {
					PosEbr.Part_next = -1
					Actual.Part_next = Actual.Part_start + Actual.Part_size
					PosEbr.Part_start = Actual.Part_start + Actual.Part_size
					dsk.Seek((Actual.Part_start), 0)
					binary.Write(bufferControl, binary.BigEndian, &Actual)
					escribirBytes(dsk, bufferControl.Bytes())
					return PosEbr
				}
			}
		}

	}
	nextpos := AllArrays[0].Part_start + AllArrays[0].Part_size
	AllArrays[0].Part_next = nextpos
	dsk.Seek((AllArrays[0].Part_start), 0)
	binary.Write(bufferControl, binary.BigEndian, &tempEbr)
	escribirBytes(dsk, bufferControl.Bytes())

	PosEbr.Part_start = nextpos
	PosEbr.Part_next = -1
	return PosEbr
}

func GetPositionFF(mbr Estructuras.MBR, size int64) int {
	mbrtemp := Estructuras.MBR{}
	posmbr := 0

	if IstheDiskEmpty(mbr) {
		Pos := int(unsafe.Sizeof(mbr)) + 1
		return Pos
	}

	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_status == '1' {
			mbrtemp.Mbr_partition[i] = mbr.Mbr_partition[i]
			posmbr += 1
		}
	}

	mbrtemp = OrdenarArray(mbrtemp)
	fin := 0

	for i := 1; i < 4; i++ {
		Actual := mbrtemp.Mbr_partition[i]
		Anterior := mbrtemp.Mbr_partition[i-1]

		if Actual.Part_status != '0' {
			fin = int(Anterior.Part_size) + int(Anterior.Part_start)
			if fin+int(size) < int(Actual.Part_start) {
				return fin
			}
		}

	}

	return int(mbrtemp.Mbr_partition[0].Part_size) + int(mbrtemp.Mbr_partition[0].Part_start)
}

func OrdenarArray(Mbr Estructuras.MBR) Estructuras.MBR {
	temp := Estructuras.MBR{}
	postemp := 0
	n := 4
	for i := 1; i < n; i++ {
		for j := n - 1; j >= i; j-- {
			if Mbr.Mbr_partition[j-1].Part_start > Mbr.Mbr_partition[j].Part_start && Mbr.Mbr_partition[j-1].Part_start != 0 && Mbr.Mbr_partition[j].Part_start != 0 {
				temp.Mbr_partition[postemp] = Mbr.Mbr_partition[j-1]
				Mbr.Mbr_partition[j-1] = Mbr.Mbr_partition[j]
				Mbr.Mbr_partition[j] = temp.Mbr_partition[postemp]
				postemp += 1
			}
		}
	}
	return Mbr
}

func IstheDiskEmpty(MBR Estructuras.MBR) bool {
	for i := 0; i < 4; i++ {
		if MBR.Mbr_partition[i].Part_status == '1' {
			return false
		}
	}

	return true
}

func ReadBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func SearchNameEBR(PartExt Estructuras.Particion, Nombre string, parameters Estructuras.ParamStruct) bool {
	fullpath := parameters.Direccion
	startpoint := PartExt.Part_start
	tempEbr := Estructuras.EBR{}

	var sizeMbr int64 = int64(unsafe.Sizeof(tempEbr))

	dsk, err := os.OpenFile(fullpath, os.O_RDWR, 0777)
	defer dsk.Close()
	dsk.Seek(startpoint, 0)

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &tempEbr)
	if err2 != nil {
		log.Fatal("binary.Read failed", err)
	}

	if tempEbr.Part_status != '1' {
		return false
	}

	for {
		if string(tempEbr.Part_name[:len(Nombre)]) == (Nombre) {
			return true
		}
		if tempEbr.Part_next == -1 {
			return false
		}
		dsk.Seek((tempEbr.Part_next), 0)

		dataControl := leerBytes(dsk, int(sizeMbr))
		bufferControl := bytes.NewBuffer(dataControl)
		errb := binary.Read(bufferControl, binary.BigEndian, &tempEbr)

		if errb != nil {
			log.Fatal("binary.Read failed", err)
		}

	}

}

func EliminarParticion(parameters Estructuras.ParamStruct) { //Falta revisar

	fullpath := parameters.Direccion

	size, _ := strconv.ParseInt(parameters.Tam, 10, 64)
	size = RealSize(size, parameters.Unit)

	if !existeDisco(fullpath) {
		fmt.Println("---No se encontro el disco buscado---")
		return
	}

	mbr := Estructuras.MBR{}
	var sizeMbr int64 = int64(unsafe.Sizeof(mbr))
	dsk, err := os.OpenFile(fullpath, os.O_RDWR, 0777)
	defer dsk.Close()
	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		return
	}

	dsk.Seek(0, 0)
	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &mbr)
	if err2 != nil {
		log.Fatal("binary.Read failed", err)
	}

	if IstheDiskEmpty(mbr) {
		fmt.Println("No se pudo encontrar la particion deseada")
	}
	extstart := Estructuras.Particion{}

	for i := 0; i < 4; i++ {
		if mbr.Mbr_partition[i].Part_type == 'e' {
			extstart = mbr.Mbr_partition[i]
		}
		if parameters.Nombre == string(mbr.Mbr_partition[i].Part_name[:len(parameters.Nombre)]) {
			for {
				fmt.Println("La particion fue encontrada, desea eliminarla? (S/N)")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				entrada := scanner.Text()

				if entrada == "S" || entrada == "s" {

					FormatPart(dsk, mbr.Mbr_partition[i])
					EmptPart := Estructuras.Particion{}
					EmptPart.Part_status = '0'
					EmptPart.Part_type = '0'
					EmptPart.Part_fit = '0'
					EmptPart.Part_size = 0
					EmptPart.Part_start = -1
					mbr.Mbr_partition[i] = EmptPart
					dsk.Seek(0, 0)
					var bufferControl bytes.Buffer
					binary.Write(&bufferControl, binary.BigEndian, &mbr)
					_, err = dsk.Write(bufferControl.Bytes())
					if err != nil {
						log.Fatal(err)
					}
					return

				} else {
					return
				}
			}
		}
	}

	if extstart.Part_status == '1' {
		EliminarPartL(extstart.Part_start, parameters.Nombre, dsk)
		return
	}

	fmt.Println("No se pudo encontrar la particion deseada")
	return
}

func EliminarPartL(startpoint int64, name string, dsk *os.File) { //Falta revisar
	tempEbr := Estructuras.EBR{}

	var sizeMbr int64 = int64(unsafe.Sizeof(tempEbr))
	dsk.Seek(0, 0)
	dsk.Seek(startpoint, 0)

	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &tempEbr)
	if err2 != nil {
		log.Fatal("binary.Read failed", err2)
	}

	if string(tempEbr.Part_name[:len(name)]) == (name) {
		for {
			fmt.Println("La particion fue encontrada, desea eliminarla? (S/N)")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			entrada := scanner.Text()

			if entrada == "S" || entrada == "s" {

				var temporal [1024]byte
				for j := 0; j < 1024; j++ {
					temporal[j] = 0
				}
				s := &temporal
				var binario bytes.Buffer
				binary.Write(&binario, binary.BigEndian, s)
				tam := tempEbr.Part_size / 1024

				for i := 0; i < int(tam); i++ {
					dsk.Seek(tempEbr.Part_start+int64(i), 0)
					_, err := dsk.Write(binario.Bytes())
					if err != nil {
						log.Fatal(err)
					}
				}
				return
			} else {
				return
			}
		}
	}

	if tempEbr.Part_next != -1 {
		EliminarPartL(tempEbr.Part_start, name, dsk)
		return
	}

	return
}

func FormatPart(dsk *os.File, Particion Estructuras.Particion) {
	var temporal [1024]byte
	for j := 0; j < 1024; j++ {
		temporal[j] = 0
	}
	s := &temporal
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	tam := Particion.Part_size / 1024

	for i := 0; i < int(tam); i++ {
		dsk.Seek(Particion.Part_start+int64(i), 0)
		_, err := dsk.Write(binario.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}

}
