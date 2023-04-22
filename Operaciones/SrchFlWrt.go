package operaciones

import (
	"Proyecto2/Estructuras"
	"os"
	"strings"
	"unsafe"
)

func WrtArchivoMkfs(SuperBlock Estructuras.Sblock, Path string, dsk *os.File, Contenido string, Replace bool) {

	InodoRoot := Estructuras.I_node{}
	InodoRoot = ReadInode(InodoRoot, SuperBlock.S_inode_start, dsk)

	pathSeparado := strings.Split(Path, "/")
	var pos int64

	for _, i := range InodoRoot.I_block {
		if (int64(i) - 1) != -1 {
			pos = (int64(i) - 1)
			WrtBlckCarpetas(pos, dsk, pathSeparado, SuperBlock, Contenido, Replace)
		}
	}

}

func WrtBlckCarpetas(pos int64, dsk *os.File, PathSeparado []string, SuperBlock Estructuras.Sblock, Contenido string, Replace bool) {
	bloqueCarpeta := Estructuras.BloqueCarpetas{}
	bloqueCarpeta = ReadBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)
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
			WrtInodoArch((int64(pos) - 1), dsk, PathSeparado, SuperBlock, Contenido, Replace)
			return
		}
	}

	return
}

func WrtInodoArch(pos int64, dsk *os.File, PathSeparado []string, SuperBlock Estructuras.Sblock, Contenido string, Replace bool) {

	InodoArchivo := ReadInode(Estructuras.I_node{}, SuperBlock.S_inode_start+SuperBlock.S_inode_size*pos, dsk)

	if string(InodoArchivo.I_type) == "1" {
		WrtFileBlocks(InodoArchivo, dsk, SuperBlock, Contenido, pos, Replace)
	} else {
		for _, i := range InodoArchivo.I_block {
			if (int64(i) - 1) != -1 {

				WrtBlckCarpetas((int64(i) - 1), dsk, PathSeparado, SuperBlock, Contenido, Replace)
			}
		}
	}

}

func WrtFileBlocks(Inodo Estructuras.I_node, dsk *os.File, SuperBlock Estructuras.Sblock, Contenido string, PosInodo int64, Replace bool) {

	chars := []byte(Contenido)
	var comp byte
	for i := 0; i < 16; i++ {
		if (int64(Inodo.I_block[i]) - 1) != -1 {
			is := int64(Inodo.I_block[i]) - 1
			pos := SuperBlock.S_block_start + (int64(unsafe.Sizeof(Estructuras.BloqueArchivos{})) * is)
			bl_archivo := ReadBloqueArchivo(Estructuras.BloqueArchivos{}, pos, dsk)
			if bl_archivo.B_content[62] == comp {
				for i := 0; i < 63; i++ {
					if bl_archivo.B_content[i] == comp || Replace {
						if len(chars) > 0 {
							bl_archivo.B_content[i] = byte(chars[0])
							chars = append(chars[1:])
						} else {
							WriteBloqueArchivo(bl_archivo, pos, dsk)
							return
						}
					}
				}
				WriteBloqueArchivo(bl_archivo, pos, dsk)
			}
		} else if (int64(Inodo.I_block[i]) - 1) == -1 {
			Contador := SuperBlock.S_blocks_count - SuperBlock.S_free_blocks_count
			SuperBlock.S_free_blocks_count = SuperBlock.S_free_blocks_count - 1
			Inodo.I_block[i] = byte(Contador + 1)
			StartInode := SuperBlock.S_inode_start + SuperBlock.S_inode_size*PosInodo
			WriteInode(Inodo, StartInode, dsk)
			is := (int64(i) - 1)
			pos := SuperBlock.S_block_start + (int64(unsafe.Sizeof(Estructuras.BloqueArchivos{})) * is)

			StartBlockPoint := SuperBlock.S_bm_inode_start - int64(unsafe.Sizeof(Estructuras.Sblock{}))
			WriteSBlock(SuperBlock, StartBlockPoint, dsk)
			bl_archivo := Estructuras.BloqueArchivos{}

			for i := 0; i < 63; i++ {
				if len(chars) > 0 {
					bl_archivo.B_content[i] = byte(chars[0])
					chars = append(chars[1:])
				} else {
					WriteBloqueArchivo(bl_archivo, pos, dsk)
					return
				}
			}
			WriteBloqueArchivo(bl_archivo, pos, dsk)
		}
	}

}
