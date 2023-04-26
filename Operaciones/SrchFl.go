package operaciones

import (
	"Proyecto2/Estructuras"
	"os"
	"strings"
	"unsafe"
)

func LeerArchivoMkfs(SuperBlock Estructuras.Sblock, Path string, dsk *os.File) string {

	InodoRoot := Estructuras.I_node{}
	InodoRoot = ReadInode(InodoRoot, SuperBlock.S_inode_start, dsk)

	pathSeparado := strings.Split(Path, "/")
	DatosArch := ""
	var pos int64

	for _, i := range InodoRoot.I_block {
		if (int64(i) - 1) != -1 {
			pos = (int64(i) - 1)
			DatosArch += ReadBlckCarpetas(pos, dsk, pathSeparado, SuperBlock)
			return DatosArch
		}
	}

	return DatosArch
}

func ReadBlckCarpetas(pos int64, dsk *os.File, PathSeparado []string, SuperBlock Estructuras.Sblock) string {
	bloqueCarpeta := ReadBloqueCarpeta(Estructuras.BloqueCarpetas{}, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)

	DatosArch := ""
	for _, i := range bloqueCarpeta.B_content {
		PathActual := [15]byte{}
		Strong := PathSeparado[0]
		copy(PathActual[:], Strong)

		if i.B_name == PathActual {
			pos := i.B_inodo
			PathSeparado = PathSeparado[1:]
			DatosArch = LeerInodoArch((int64(pos) - 1), dsk, PathSeparado, SuperBlock)
			return DatosArch
		}
	}

	return DatosArch
}

func LeerInodoArch(pos int64, dsk *os.File, PathSeparado []string, SuperBlock Estructuras.Sblock) string {
	DatosArch := ""

	InodoArchivo := ReadInode(Estructuras.I_node{}, SuperBlock.S_inode_start+SuperBlock.S_inode_size*pos, dsk)

	if string(InodoArchivo.I_type) == "1" {
		DatosArch += ReadFileBlocks(InodoArchivo, dsk, SuperBlock)
	} else {
		for _, i := range InodoArchivo.I_block {
			if (int64(i) - 1) != -1 {

				DatosArch += ReadBlckCarpetas((int64(i) - 1), dsk, PathSeparado, SuperBlock)
			}
		}
	}

	return DatosArch

}

func ReadFileBlocks(Inodo Estructuras.I_node, dsk *os.File, SuperBlock Estructuras.Sblock) string {
	DatosArch := ""

	for _, i := range Inodo.I_block {
		if (int64(i) - 1) != -1 {
			is := (int64(i) - 1)
			pos := SuperBlock.S_block_start + (int64(unsafe.Sizeof(Estructuras.BloqueArchivos{})) * is)
			bl_archivo := ReadBloqueArchivo(Estructuras.BloqueArchivos{}, pos, dsk)
			DatosArch += string((bl_archivo.B_content[:]))
		}
	}

	return DatosArch
}
