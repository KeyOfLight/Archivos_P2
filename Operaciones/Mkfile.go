package operaciones

import (
	"Proyecto2/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Mkfile(parameters Estructuras.ParamStruct) {
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

	size, _ := strconv.ParseInt(parameters.Size, 10, 64)
	base := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var contenido string = ""
	pivote := 0
	for i := 0; i < int(size); i++ {
		contenido += string(fmt.Sprint(base[pivote]))
		if pivote == 9 {
			pivote = 0
		}
		pivote++
	}

	pathSeparado := strings.Split(VirtualPath, "/")
	if pathSeparado[0] == "" {
		pathSeparado = append(pathSeparado[1:])
	}

	if parameters.Tam != "" {
		Datos, err := ioutil.ReadFile(parameters.Tam)

		if err != nil {
			log.Fatal(err)
		}
		contenido = ""
		datosComoString := string(Datos)
		SeparadoLineas := strings.Split(datosComoString, "\n")
		for _, linea := range SeparadoLineas {
			contenido += linea
		}

		size = int64(unsafe.Sizeof(contenido))
	}
	r := false
	if parameters.Pwd != "" {
		r = true
	}

	pos := SrchInodo(0, dsk, pathSeparado, StartPoint, r, size)

	if pos == -1 {
		fmt.Println("No se pudo encontrar la ubicacion de archivo")
		return
	}

	SuperBlock := ReadSBlock(Estructuras.Sblock{}, StartPoint, dsk)
	WrtArchivoMkfs(SuperBlock, VirtualPath, dsk, contenido, true)
}

func SrchBlckCarpetas(pos int64, dsk *os.File, PathSeparado []string, SuperBlockSp int64, CreateFull bool, Size int64) (Pos int64, Econtrado bool) {
	bloqueCarpeta := Estructuras.BloqueCarpetas{}
	SuperBlock := ReadSBlock(Estructuras.Sblock{}, SuperBlockSp, dsk)
	bloqueCarpeta = ReadBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)
	Coincidencia_Path := false

	for _, i := range bloqueCarpeta.B_content {
		PathActual := [15]byte{}
		if len(PathSeparado) == 0 {
			return 0, true
		}

		Strong := PathSeparado[0]
		copy(PathActual[:], Strong)

		if i.B_name == PathActual {
			pos := i.B_inodo
			PathSeparado = PathSeparado[1:]
			Coincidencia_Path = true
			return SrchInodo((int64(pos) - 1), dsk, PathSeparado, SuperBlockSp, true, Size), Coincidencia_Path
		}
	}

	return -1, Coincidencia_Path
}

func SrchInodo(pos int64, dsk *os.File, PathSeparado []string, SuperBlockSp int64, CreateFull bool, Size int64) int64 {
	SuperBlock := ReadSBlock(Estructuras.Sblock{}, SuperBlockSp, dsk)
	StartInode := SuperBlock.S_inode_start + SuperBlock.S_inode_size*pos
	InodoArchivo := ReadInode(Estructuras.I_node{}, StartInode, dsk)
	Encontrado := false
	ReturnedPos := int64(0)
	if len(PathSeparado) == 0 {
		return pos
	} else {
		for _, i := range InodoArchivo.I_block {
			if (int64(i) - 1) != -1 {
				ReturnedPos, Encontrado = SrchBlckCarpetas((int64(i) - 1), dsk, PathSeparado, SuperBlockSp, CreateFull, Size)
			}
		}
		if Encontrado {
			return ReturnedPos
		}

		if !CreateFull {
			return -1
		}

		var Retornado = int64(-1)
		for i := 0; i < 16; i++ {
			if int64(InodoArchivo.I_block[i])-1 != -1 {
				Retornado, _ = findEmptySpaceInCarpet(int64(InodoArchivo.I_block[i])-1, dsk, PathSeparado, SuperBlockSp, Size, false)
				Next := PathSeparado[1]
				if strings.Contains(Next, ".") {
					MkInodeCarpeta(Retornado, dsk, Size, SuperBlock, true)
				} else {
					MkInodeCarpeta(Retornado, dsk, Size, SuperBlock, false)
				}

				SrchBlckCarpetas(int64(InodoArchivo.I_block[i])-1, dsk, PathSeparado, SuperBlockSp, CreateFull, Size)
			}

			if Retornado != -1 {
				return Retornado
			}

			if int64(InodoArchivo.I_block[i])-1 == -1 {
				var posBlock int64
				Retornado, posBlock = findEmptySpaceInCarpet(int64(InodoArchivo.I_block[i])-1, dsk, PathSeparado, SuperBlockSp, Size, true)
				Next := PathSeparado[1]
				InodoArchivo.I_block[i] = byte(posBlock + 1)
				WriteInode(InodoArchivo, SuperBlock.S_inode_start+SuperBlock.S_inode_size*pos, dsk)

				if strings.Contains(Next, ".") {
					MkInodeCarpeta(Retornado, dsk, Size, SuperBlock, true)
				} else {
					MkInodeCarpeta(Retornado, dsk, Size, SuperBlock, false)
					SrchBlckCarpetas(int64(InodoArchivo.I_block[i])-1, dsk, PathSeparado, SuperBlockSp, CreateFull, Size)
				}

			}

		}

	}

	return -1
}

func findEmptySpaceInCarpet(pos int64, dsk *os.File, PathSeparado []string, SuperBlockSp int64, Size int64, Inicial bool) (PosNueva int64, PosCarpeta int64) {
	SuperBlock := ReadSBlock(Estructuras.Sblock{}, SuperBlockSp, dsk)
	if Inicial {
		bloqueCarpeta := Estructuras.BloqueCarpetas{}

		NewPosblock := SuperBlock.S_first_blo + 1
		SuperBlock.S_first_blo += 1
		SuperBlock.S_free_blocks_count = SuperBlock.S_free_blocks_count - 1

		Strong := PathSeparado[0]
		NewPos := SuperBlock.S_firts_ino + 1
		SuperBlock.S_firts_ino += 1

		SuperBlock.S_free_inodes_count = SuperBlock.S_free_inodes_count - 1
		StartBlockPoint := SuperBlock.S_bm_inode_start - int64(unsafe.Sizeof(Estructuras.Sblock{}))

		copy(bloqueCarpeta.B_content[2].B_name[:], []byte(Strong))
		bloqueCarpeta.B_content[2].B_inodo = byte(NewPos + 1)

		bloqueCarpeta.B_content[0].B_inodo = 1
		copy(bloqueCarpeta.B_content[0].B_name[:], []byte("."))

		bloqueCarpeta.B_content[1].B_inodo = 1
		copy(bloqueCarpeta.B_content[1].B_name[:], []byte(".."))

		bloqueCarpeta.B_content[3].B_inodo = 0
		copy(bloqueCarpeta.B_content[3].B_name[:], []byte(""))

		WriteSBlock(SuperBlock, StartBlockPoint, dsk)
		WriteBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*NewPosblock, dsk)
		AddBlockBm(dsk, SuperBlock)

		return NewPos, NewPosblock
	}

	bloqueCarpeta := Estructuras.BloqueCarpetas{}
	bloqueCarpeta = ReadBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)

	for i := 0; i < 4; i++ {
		if int64(bloqueCarpeta.B_content[i].B_inodo)-1 == -1 {
			Strong := PathSeparado[0]
			NewPos := SuperBlock.S_firts_ino + 1
			SuperBlock.S_firts_ino += 1
			copy(bloqueCarpeta.B_content[i].B_name[:], []byte(Strong))
			bloqueCarpeta.B_content[i].B_inodo = byte(NewPos + 1)

			SuperBlock.S_free_inodes_count = SuperBlock.S_free_inodes_count - 1
			StartBlockPoint := SuperBlock.S_bm_inode_start - int64(unsafe.Sizeof(Estructuras.Sblock{}))

			WriteSBlock(SuperBlock, StartBlockPoint, dsk)
			WriteBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)
			AddBlockBm(dsk, SuperBlock)

			return NewPos, 0

		}
	}

	return -1, -1

}

func MkInodeCarpeta(posInode int64, dsk *os.File, size int64, SuperBlock Estructuras.Sblock, File bool) { //Verificar si aun es necesario
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()

	dt := time.Now()
	time := dt.String()

	Uid, _ := strconv.ParseInt(string(Uss.UID), 10, 64)
	Gid, _ := strconv.ParseInt(string(Uss.GUID), 10, 64)

	NewinodoArchivo := Estructuras.I_node{}
	NewinodoArchivo.I_uid = Uid
	NewinodoArchivo.I_gid = Gid
	NewinodoArchivo.I_size = size
	NewinodoArchivo.I_perm = 664
	copy(NewinodoArchivo.I_atime[:], []byte(time))
	copy(NewinodoArchivo.I_ctime[:], []byte(time))
	copy(NewinodoArchivo.I_mtime[:], []byte(time))

	for i := 0; i < 16; i++ {
		NewinodoArchivo.I_block[i] = 0
	}

	if File {
		NewinodoArchivo.I_type = '1'

	} else {
		NewinodoArchivo.I_type = '0'
	}
	StartInode := SuperBlock.S_inode_start + SuperBlock.S_inode_size*posInode
	WriteInode(NewinodoArchivo, StartInode, dsk)
	AddInodeBm(dsk, SuperBlock)

}

func AddInodeBm(dsk *os.File, SuperBlock Estructuras.Sblock) {
	var bit byte
	for i := SuperBlock.S_bm_inode_start; i < (SuperBlock.S_bm_block_start); i++ {
		dsk.Seek(i, 0)
		dataControl := ReadBytes(dsk, 1)
		bufferControl := bytes.NewBuffer(dataControl)
		err2 := binary.Read(bufferControl, binary.BigEndian, &bit)
		if err2 != nil {
			log.Fatal("binary.Read failed", err2)
		}

		if bit == '0' {
			bit = 1
			var bufferControl bytes.Buffer
			binary.Write(&bufferControl, binary.BigEndian, &bit)
			_, err := dsk.Write(bufferControl.Bytes())

			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	return
}

func AddBlockBm(dsk *os.File, SuperBlock Estructuras.Sblock) {
	var bit byte
	for i := SuperBlock.S_bm_block_start; i < (SuperBlock.S_inode_start); i++ {
		dsk.Seek(i, 0)
		dataControl := ReadBytes(dsk, 1)
		bufferControl := bytes.NewBuffer(dataControl)
		err2 := binary.Read(bufferControl, binary.BigEndian, &bit)
		if err2 != nil {
			log.Fatal("binary.Read failed", err2)
		}

		if bit == '0' {
			bit = 1
			var bufferControl bytes.Buffer
			binary.Write(&bufferControl, binary.BigEndian, &bit)
			_, err := dsk.Write(bufferControl.Bytes())

			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	return

}
