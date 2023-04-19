package operaciones

import (
	"Proyecto2/Estructuras"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Mkfile(parameters Estructuras.ParamStruct) {
	Path := parameters.Direccion
	r := false

	if parameters.Pwd != "" {
		r = true
	}

	if r {
		makealldirs(parameters)
	} else {
		file, err2 := os.Create(Path)
		defer file.Close()
		if err2 != nil {
			fmt.Println("No se pudo crear el archivo deseado")
			fmt.Println(err2.Error())
			return
		}
	}
}

func makealldirs(parameters Estructuras.ParamStruct) {
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
	size, _ := strconv.ParseInt(parameters.Size, 10, 64)
	pos := SrchInodo(0, dsk, pathSeparado, SuperBlock, true, size)

	fmt.Println(pos)

}

func SrchBlckCarpetas(pos int64, dsk *os.File, PathSeparado []string, SuperBlock Estructuras.Sblock, CreateFull bool, Size int64) (Pos int64, Econtrado bool) {
	bloqueCarpeta := Estructuras.BloqueCarpetas{}
	bloqueCarpeta = ReadBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)
	Coincidencia_Path := false

	if PathSeparado[0] == "" {
		PathSeparado = append(PathSeparado[1:])
	}

	for _, i := range bloqueCarpeta.B_content {
		PathActual := [15]byte{}

		Strong := PathSeparado[0]
		copy(PathActual[:], Strong)

		if i.B_name == PathActual {
			pos := i.B_inodo
			PathSeparado = PathSeparado[1:]
			Coincidencia_Path = true
			return SrchInodo((int64(pos) - 1), dsk, PathSeparado, SuperBlock, true, Size), Coincidencia_Path
		}
	}

	return -1
}

func SrchInodo(pos int64, dsk *os.File, PathSeparado []string, SuperBlock Estructuras.Sblock, CreateFull bool, Size int64) int64 {
	InodoArchivo := ReadInode(Estructuras.I_node{}, SuperBlock.S_inode_start+SuperBlock.S_inode_size*pos, dsk)
	Encontrado := false
	ReturnedPos := int64(0)
	if len(PathSeparado) == 0 {
		return pos
	} else {
		for _, i := range InodoArchivo.I_block {
			if (int64(i) - 1) != -1 { //TIene que buscar en todos lados antes de crear, cambiar casi todo basicamente
				ReturnedPos, Encontrado = SrchBlckCarpetas((int64(i) - 1), dsk, PathSeparado, SuperBlock, CreateFull, Size)
			}
		}
		if Encontrado {
			return ReturnedPos
		}

		if CreateFull {
			for i := 0; i < 16; i++ {
				if int64(InodoArchivo.I_block[i])-1 != -1 {

				}
			}
		}

	}

	return -1
}

func findEmptySpaceInCarpet(pos int64, dsk *os.File, PathSeparado []string, SuperBlock Estructuras.Sblock, Size int64) {

	bloqueCarpeta := Estructuras.BloqueCarpetas{}
	bloqueCarpeta = ReadBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)

	for i := 0; i < 4; i++ {
		if int64(bloqueCarpeta.B_content[i].B_inodo)-1 == -1 {
			Strong := PathSeparado[0]
			NewPos := SuperBlock.S_inodes_count - SuperBlock.S_free_inodes_count + 1
			copy(bloqueCarpeta.B_content[i].B_name[:], []byte(Strong))
			bloqueCarpeta.B_content[i].B_inodo = byte(NewPos)

			SuperBlock.S_free_inodes_count = SuperBlock.S_free_inodes_count - 1
			StartBlockPoint := SuperBlock.S_bm_inode_start - int64(unsafe.Sizeof(Estructuras.Sblock{}))

			WriteSBlock(SuperBlock, StartBlockPoint, dsk)
			WriteBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)

			Next := PathSeparado[1]

			if strings.Contains(Next, ".") {
				return NewPos
			} else {
				PosInode := MkInodeCarpeta(NewPos, dsk, Size, SuperBlock)
				return SrchInodo((PosInode - 1), dsk, PathSeparado, SuperBlock, true, Size)
			}

		}
	}

}

func mkBlockCarpeta(dsk *os.File, size int64, SuperBlock Estructuras.Sblock) int64 {

}

func MkInodeCarpeta(posInode int64, dsk *os.File, size int64, SuperBlock Estructuras.Sblock) int64 {
	var Uss = Estructuras.User{}
	Uss = (&Uss).Getusser()

	NewPos := SuperBlock.S_blocks_count - SuperBlock.S_blocks_count + 1
	SuperBlock.S_free_blocks_count = SuperBlock.S_free_blocks_count - 1
	StartBlockPoint := SuperBlock.S_bm_inode_start - int64(unsafe.Sizeof(Estructuras.Sblock{}))
	WriteSBlock(SuperBlock, StartBlockPoint, dsk)

	dt := time.Now()
	time := dt.String()

	Uid, _ := strconv.ParseInt(string(Uss.UID), 10, 64)
	Gid, _ := strconv.ParseInt(string(Uss.Grupo[:]), 10, 64)

	NewinodoArchivo := Estructuras.I_node{}
	NewinodoArchivo.I_uid = Uid
	NewinodoArchivo.I_gid = Gid
	NewinodoArchivo.I_size = size
	NewinodoArchivo.I_perm = 664
	NewinodoArchivo.I_type = '0'
	copy(NewinodoArchivo.I_atime[:], []byte(time))
	copy(NewinodoArchivo.I_ctime[:], []byte(time))
	copy(NewinodoArchivo.I_mtime[:], []byte(time))

	NewinodoArchivo.I_block[0] = byte(NewPos)
	for i := 1; i < 16; i++ {
		NewinodoArchivo.I_block[i] = 0
	}
	StartInode := SuperBlock.S_inode_start + SuperBlock.S_inode_size*posInode
	WriteInode(NewinodoArchivo, StartInode, dsk)

	BloqueCarpeta1 := Estructuras.BloqueCarpetas{}

	BloqueCarpeta1.B_content[0].B_inodo = 1
	copy(BloqueCarpeta1.B_content[0].B_name[:], []byte("."))

	BloqueCarpeta1.B_content[1].B_inodo = 1
	copy(BloqueCarpeta1.B_content[1].B_name[:], []byte(".."))

	BloqueCarpeta1.B_content[2].B_inodo = 2
	copy(BloqueCarpeta1.B_content[2].B_name[:], []byte("users.txt"))

	BloqueCarpeta1.B_content[3].B_inodo = 0
	copy(BloqueCarpeta1.B_content[3].B_name[:], []byte(""))

	return posInode
}

func MkFileBlocks(pos int64, dsk *os.File, size int64, SuperBlock Estructuras.Sblock, User Estructuras.User) {
	InodoArchivo := ReadInode(Estructuras.I_node{}, SuperBlock.S_inode_start+SuperBlock.S_inode_size*pos, dsk)
	NewPos := SuperBlock.S_blocks_count - SuperBlock.S_free_blocks_count + 1
	SuperBlock.S_free_blocks_count = SuperBlock.S_free_blocks_count - 1
	StartBlockPoint := SuperBlock.S_bm_inode_start - int64(unsafe.Sizeof(Estructuras.Sblock{}))
	WriteSBlock(SuperBlock, StartBlockPoint, dsk)

	for i := 0; i < len(InodoArchivo.I_block); i++ {
		if int64(InodoArchivo.I_block[i])-1 == -1 {
			InodoArchivo.I_block[i] = byte(NewPos)
			break
		}
	}

	dt := time.Now()
	time := dt.String()
	copy(InodoArchivo.I_atime[:], []byte(time))
	copy(InodoArchivo.I_ctime[:], []byte(time))
	copy(InodoArchivo.I_mtime[:], []byte(time))
	StartInode := SuperBlock.S_inode_start + SuperBlock.S_inode_size*pos
	WriteInode(InodoArchivo, StartInode, dsk)

	Uid, _ := strconv.ParseInt(string(User.UID), 10, 64)
	Gid, _ := strconv.ParseInt(string(User.Grupo[:]), 10, 64)

	NewinodoArchivo := Estructuras.I_node{}
	NewinodoArchivo.I_uid = Uid
	NewinodoArchivo.I_gid = Gid
	NewinodoArchivo.I_size = size
	NewinodoArchivo.I_perm = 664
	NewinodoArchivo.I_type = '1'
	copy(NewinodoArchivo.I_atime[:], []byte(time))
	copy(NewinodoArchivo.I_ctime[:], []byte(time))
	copy(NewinodoArchivo.I_mtime[:], []byte(time))

}
