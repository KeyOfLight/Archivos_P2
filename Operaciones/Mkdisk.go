package operaciones

import (
	"Proyecto2/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func MakeDisk(parameters Estructuras.ParamStruct) {
	fullpath := parameters.Direccion

	if existeDisco(fullpath) {
		fmt.Println("*** El disco ya existe ***")
		return
	}

	fsize, _ := strconv.ParseInt(parameters.Tam, 10, 64)
	size, _ := strconv.ParseInt(parameters.Tam, 10, 64)
	size = RealSize(size, parameters.Unit)

	mbr := Estructuras.MBR{}
	mbr.Mbr_tamano = size

	mbr.Dsk_fit = getFit(parameters.Fit)[0]
	mbr.Mbr_dsk_signature = int64(rand.Int())

	dt := time.Now()
	time := dt.String()

	for t := 0; t < 19; t++ {
		mbr.Mbr_fecha_creacion[t] = time[t]
	}

	for p := 0; p < 4; p++ {
		mbr.Mbr_partition[p].Part_status = '0'
		mbr.Mbr_partition[p].Part_type = '0'
		mbr.Mbr_partition[p].Part_fit = '0'
		mbr.Mbr_partition[p].Part_size = 0
		mbr.Mbr_partition[p].Part_start = -1
	}

	crearDirectorio(fullpath)

	file, err2 := os.Create(fullpath)
	defer file.Close()
	if err2 != nil {
		fmt.Println("No se pudo crear el archivo deseado")
		fmt.Println(err2.Error())
		return
	}

	var temporal [1024]byte
	for i := 0; i < 1024; i++ {
		temporal[i] = 0
	}
	s := &temporal
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	tam := 0
	if parameters.Unit == "k" {
		tam = int(fsize)
	} else if parameters.Unit == "m" {
		tam = int(fsize) * 1024
	} else {
		tam = int(fsize) * 1024
	}

	for i := 0; i < tam; i++ {
		_, err := file.Write(binario.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}

	file.Seek(0, 0)
	var bufferControl bytes.Buffer
	binary.Write(&bufferControl, binary.BigEndian, &mbr)
	_, err := file.Write(bufferControl.Bytes())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Disco creado correctamente")
	fmt.Println("Se agrego el mbr de manera correcta")
	return

}

func crearDirectorio(path string) {
	directorio := PathClanner(path)
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.MkdirAll(directorio, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func PathClanner(fullpath string) string {

	separado := strings.Split(fullpath, "/")
	separado = separado[:len(separado)-1]
	var regresar string
	for _, linea := range separado {
		if linea != "" {
			regresar += "/" + linea
		}
	}

	return regresar
}

func getFit(Fit string) []byte {

	if Fit == "FF" {
		bi := []byte("f")
		return bi
	} else if Fit == "WF" {
		bi := []byte("w")
		return bi
	} else if Fit == "BF" {
		bi := []byte("b")
		return bi
	}

	bi := []byte("f")
	return bi

}

func RealSize(tama int64, Unit string) int64 {

	if Unit == "k" {
		tama = tama * 1024
	} else if Unit == "m" {
		tama = tama * 1024 * 1024
	} else {
		tama = tama * 1024 * 1024
	}
	var newsize = tama
	return newsize
}
func existeDisco(path string) bool {
	_, err := ioutil.ReadFile(path)

	if err != nil {
		return false
	}

	return true
}
