package operaciones

import (
	"Proyecto2/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"
)

func VerFit(CurrentFit string) int64 {

	if CurrentFit == "ff" || CurrentFit == "" {
		return 1
	} else if CurrentFit == "wf" {
		return 2
	} else if CurrentFit == "bf" {
		return 3
	}

	return 1
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func WriteSBlock(data Estructuras.Sblock, StartPoint int64, dsk *os.File) {

	dsk.Seek(StartPoint, 0)
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &data)
	_, err := dsk.Write(bufferControl.Bytes())

	if err != nil {
		log.Fatal(err)
	}
}

func ReadSBlock(data Estructuras.Sblock, StartPoint int64, dsk *os.File) Estructuras.Sblock {

	var sizeMbr int64 = int64(unsafe.Sizeof(data))
	dsk.Seek(StartPoint, 0)
	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &data)
	if err2 != nil {
		log.Fatal("binary.Read failed", err2)
	}

	return data
}

func WriteInode(data Estructuras.I_node, StartPoint int64, dsk *os.File) {

	dsk.Seek(StartPoint, 0)
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &data)
	_, err := dsk.Write(bufferControl.Bytes())

	if err != nil {
		log.Fatal(err)
	}
}

func ReadInode(data Estructuras.I_node, StartPoint int64, dsk *os.File) Estructuras.I_node {

	var sizeMbr int64 = int64(unsafe.Sizeof(data))
	dsk.Seek(StartPoint, 0)
	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &data)
	if err2 != nil {
		log.Fatal("binary.Read failed", err2)
	}

	dt := time.Now()
	time := dt.String()
	copy(data.I_atime[:], []byte(time))
	WriteInode(data, StartPoint, dsk)

	return data
}

func WriteBloqueCarpeta(data Estructuras.BloqueCarpetas, StartPoint int64, dsk *os.File) {

	dsk.Seek(StartPoint, 0)
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &data)
	_, err := dsk.Write(bufferControl.Bytes())

	if err != nil {
		log.Fatal(err)
	}
}

func ReadBloqueCarpeta(data Estructuras.BloqueCarpetas, StartPoint int64, dsk *os.File) Estructuras.BloqueCarpetas {

	var sizeMbr int64 = int64(unsafe.Sizeof(data))
	dsk.Seek(StartPoint, 0)
	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &data)
	if err2 != nil {
		log.Fatal("binary.Read failed", err2)
	}

	return data
}

func WriteBloqueArchivo(data Estructuras.BloqueArchivos, StartPoint int64, dsk *os.File) {

	dsk.Seek(StartPoint, 0)
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &data)
	_, err := dsk.Write(bufferControl.Bytes())

	if err != nil {
		log.Fatal(err)
	}
}

func ReadBloqueArchivo(data Estructuras.BloqueArchivos, StartPoint int64, dsk *os.File) Estructuras.BloqueArchivos {

	var sizeMbr int64 = int64(unsafe.Sizeof(data))
	dsk.Seek(StartPoint, 0)
	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &data)
	if err2 != nil {
		log.Fatal("binary.Read failed", err2)
	}

	return data
}

func OpenMBR(fullpath string) Estructuras.MBR {
	mbr := Estructuras.MBR{}

	var sizeMbr int64 = int64(unsafe.Sizeof(mbr))

	dsk, err := os.OpenFile(fullpath, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	err2 := binary.Read(bufferControl, binary.BigEndian, &mbr)
	if err2 != nil {
		log.Fatal("binary.Read failed", err)
	}

	return mbr
}

func ReadEbr(startpoint int64, dsk *os.File) Estructuras.EBR {

	dsk.Seek((startpoint), 0)
	tempEbr := Estructuras.EBR{}
	var sizeMbr int64 = int64(unsafe.Sizeof(tempEbr))

	dataControl := ReadBytes(dsk, int(sizeMbr))
	bufferControl := bytes.NewBuffer(dataControl)
	errb := binary.Read(bufferControl, binary.BigEndian, &tempEbr)
	if errb != nil {
		log.Fatal("binary.Read failed", errb)
	}

	return tempEbr
}
