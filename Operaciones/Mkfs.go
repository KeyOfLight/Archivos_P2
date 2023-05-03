package operaciones

import (
	"Proyecto2/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"time"
	"unsafe"
)

func MakeFs(parameters Estructuras.ParamStruct) {

	var Listado = Estructuras.ListaMontados{}
	Listado = (&Listado).GetLista()
	for _, i := range Listado.Montado {
		if parameters.Nombre == i.Id {
			if parameters.Fs == "ext3" {

			} else {
				Mkfs2(i)
				return
			}
		}
	}
}

func Mkfs2(Particion Estructuras.PartMounted) {
	path := Particion.Path
	dt := time.Now()
	time := dt.String()
	StartPoint := Particion.StartPoint
	nulo := '0'

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	n := (float64(Particion.Size) - float64(unsafe.Sizeof(Estructuras.Sblock{}))) / (4 + float64(unsafe.Sizeof(Estructuras.I_node{})) + 3*float64(unsafe.Sizeof(Estructuras.BloqueArchivos{})))

	if n < 1 {
		fmt.Println("La particion es demasiado peque;a")
	}

	Estructs_Num := math.Floor(n)
	num_block := 3 * Estructs_Num

	SupBlock := Estructuras.Sblock{}

	SupBlock.S_filesystem_type = 2
	SupBlock.S_inodes_count = int64(Estructs_Num)
	SupBlock.S_blocks_count = int64(num_block)
	SupBlock.S_free_inodes_count = int64(Estructs_Num - 2)
	SupBlock.S_free_blocks_count = int64(num_block - 2)
	for i := 0; i < 19; i++ {
		SupBlock.S_mtime[i] = time[i]
	}

	SupBlock.S_mnt_count = 0
	SupBlock.S_magic = 0xEF53
	SupBlock.S_inode_size = int64(unsafe.Sizeof(Estructuras.I_node{}))
	SupBlock.S_block_size = 64
	SupBlock.S_firts_ino = 2
	SupBlock.S_first_blo = 2

	PosSupB := StartPoint
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &nulo)

	if Particion.Tipo == 'e' {
		PosSupB = StartPoint + int64(unsafe.Sizeof(Estructuras.EBR{}))
		SupBlock.S_block_start = StartPoint + int64(unsafe.Sizeof(Estructuras.EBR{})) + int64(unsafe.Sizeof(Estructuras.Sblock{})) + int64(Estructs_Num) + int64(num_block) + (int64(unsafe.Sizeof(Estructuras.I_node{})) * int64(Estructs_Num))
		SupBlock.S_bm_inode_start = StartPoint + int64(unsafe.Sizeof(Estructuras.Sblock{})) + int64(unsafe.Sizeof(Estructuras.EBR{}))
		SupBlock.S_bm_block_start = StartPoint + int64(unsafe.Sizeof(Estructuras.Sblock{})) + int64(Estructs_Num) + int64(unsafe.Sizeof(Estructuras.EBR{}))
		SupBlock.S_inode_start = StartPoint + int64(unsafe.Sizeof(Estructuras.Sblock{})) + int64(Estructs_Num) + int64(num_block) + int64(unsafe.Sizeof(Estructuras.EBR{}))

		dsk.Seek(PosSupB, 0)
	} else {

		SupBlock.S_block_start = StartPoint + int64(unsafe.Sizeof(Estructuras.Sblock{})) + int64(Estructs_Num) + int64(num_block) + (int64(unsafe.Sizeof(Estructuras.I_node{})) * int64(Estructs_Num))
		SupBlock.S_bm_inode_start = StartPoint + int64(unsafe.Sizeof(Estructuras.Sblock{}))
		SupBlock.S_bm_block_start = StartPoint + int64(unsafe.Sizeof(Estructuras.Sblock{})) + int64(Estructs_Num)
		SupBlock.S_inode_start = StartPoint + int64(unsafe.Sizeof(Estructuras.Sblock{})) + int64(Estructs_Num) + int64(num_block)

		dsk.Seek(StartPoint, 0)
	}

	WriteSBlock(SupBlock, StartPoint, dsk)

	inodos := make([]byte, int64(n))
	bloques := make([]byte, 3*int64(n))

	for i := 0; i < int(n); i++ {
		inodos[i] = '0'
	}
	inodos[0] = '1'
	inodos[1] = '1'

	dsk.Seek(PosSupB+int64(unsafe.Sizeof(Estructuras.Sblock{})), 0)
	binary.Write(&binario, binary.BigEndian, &inodos)
	_, err = dsk.Write(binario.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < int(n)*3; i++ {
		bloques[i] = '0'
	}
	bloques[0] = '1'
	bloques[1] = '1'

	dsk.Seek(PosSupB+int64(unsafe.Sizeof(Estructuras.Sblock{}))+int64(unsafe.Sizeof(Estructuras.I_node{})), 0)
	binary.Write(&binario, binary.BigEndian, &bloques)
	_, err = dsk.Write(binario.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	inodoRoot := Estructuras.I_node{}

	inodoRoot.I_uid = 1
	inodoRoot.I_gid = 1
	inodoRoot.I_size = 0

	for i := 0; i < 19; i++ {
		inodoRoot.I_atime[i] = time[i]
		inodoRoot.I_ctime[i] = time[i]
		inodoRoot.I_mtime[i] = time[i]
	}

	inodoRoot.I_block[0] = 1
	inodoRoot.I_type = '0'
	inodoRoot.I_perm = 664

	for i := 1; i < 16; i++ {
		inodoRoot.I_block[i] = 0
	}

	BloqueCarpeta1 := Estructuras.BloqueCarpetas{}

	BloqueCarpeta1.B_content[0].B_inodo = 1
	copy(BloqueCarpeta1.B_content[0].B_name[:], []byte("."))

	BloqueCarpeta1.B_content[1].B_inodo = 1
	copy(BloqueCarpeta1.B_content[1].B_name[:], []byte(".."))

	BloqueCarpeta1.B_content[2].B_inodo = 2
	copy(BloqueCarpeta1.B_content[2].B_name[:], []byte("users.txt"))

	BloqueCarpeta1.B_content[3].B_inodo = 0
	copy(BloqueCarpeta1.B_content[3].B_name[:], []byte(""))

	contenido := "1,G,root\n1,U,root,root,123\n"

	inodoArchivo := Estructuras.I_node{}
	inodoArchivo.I_uid = 1
	inodoArchivo.I_gid = 1
	inodoArchivo.I_size = int64(unsafe.Sizeof(contenido))
	inodoArchivo.I_perm = 664
	inodoArchivo.I_type = '1'
	copy(inodoArchivo.I_atime[:], []byte(time))
	copy(inodoArchivo.I_ctime[:], []byte(time))
	copy(inodoArchivo.I_mtime[:], []byte(time))

	inodoArchivo.I_block[0] = 2
	for i := 1; i < 16; i++ {
		inodoArchivo.I_block[i] = 0
	}

	bl_archivo := Estructuras.BloqueArchivos{}
	for i := 0; i < len(bl_archivo.B_content); i++ {
		bl_archivo.B_content[0] = 0
	}
	copy(bl_archivo.B_content[:], contenido)

	//Escribir Inodo Root
	WriteInode(inodoRoot, SupBlock.S_inode_start, dsk)

	//Escribir Carpeta Root
	WriteBloqueCarpeta(BloqueCarpeta1, SupBlock.S_block_start, dsk)

	//Escribir Inodo Archivo
	WriteInode(inodoArchivo, SupBlock.S_inode_start+int64(unsafe.Sizeof(Estructuras.I_node{})), dsk)

	//Escribir Bloque del Archivo
	pos := SupBlock.S_block_start + int64(unsafe.Sizeof(Estructuras.BloqueArchivos{}))
	WriteBloqueArchivo(bl_archivo, pos, dsk)

	fmt.Println(" Particion formateada en EXT2")
}
