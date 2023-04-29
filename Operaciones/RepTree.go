package operaciones

import (
	"Proyecto2/Estructuras"
	"fmt"
	"os"
	"unsafe"
)

func MkTree(SuperBlock Estructuras.Sblock, dsk *os.File) string {

	InodoRoot := Estructuras.I_node{}
	InodoRoot = ReadInode(InodoRoot, SuperBlock.S_inode_start, dsk)
	DatosArch := "I0[label=<\n<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"1\" CELLPADDING=\"4\">\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"purple\" COLSPAN=\"2\">Inodo Root</TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_uid</TD><TD>" + fmt.Sprint(InodoRoot.I_uid) + " </TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_gid</TD><TD>" + fmt.Sprint(InodoRoot.I_gid) + " </TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_size</TD><TD>" + fmt.Sprint(InodoRoot.I_size) + " </TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_atime</TD><TD>" + string(InodoRoot.I_atime[:]) + " </TD>\n</TR>\n"

	var pos int64

	for i := 0; i < 16; i++ {
		act := InodoRoot.I_block[i]
		if (int64(act) - 1) != -1 {
			DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">Ap" + fmt.Sprint(i) + "</TD><TD>" + fmt.Sprint(act) + " </TD>\n</TR>\n"
		}
	}

	DatosArch += "</TABLE>>];\n"

	for _, i := range InodoRoot.I_block {
		if (int64(i) - 1) != -1 {
			pos = (int64(i) - 1)
			DatosArch += "I0->" + "B" + fmt.Sprint(pos) + "\n"
			DatosArch += TreeBlckCarpetas(pos, dsk, SuperBlock)
		}
	}

	return DatosArch
}

func TreeBlckCarpetas(pos int64, dsk *os.File, SuperBlock Estructuras.Sblock) string {
	bloqueCarpeta := Estructuras.BloqueCarpetas{}
	bloqueCarpeta = ReadBloqueCarpeta(bloqueCarpeta, SuperBlock.S_block_start+SuperBlock.S_block_size*pos, dsk)

	DatosArch := "B" + fmt.Sprint(pos) + "[label=<\n<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"1\" CELLPADDING=\"4\">\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"purple\" COLSPAN=\"2\">Bloque Carpeta</TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">Nombre</TD><TD>Pos</TD>\n</TR>\n"
	for i := 0; i < 4; i++ {
		act := bloqueCarpeta.B_content[i]
		DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">" + remove0s(act.B_name[:]) + "</TD><TD>" + fmt.Sprint(act.B_inodo) + " </TD>\n</TR>\n"
	}
	DatosArch += "</TABLE>>];\n"

	for _, i := range bloqueCarpeta.B_content {
		pos2 := i.B_inodo
		if (int(pos2)-1) != 0 && (int(pos2)-1) != -1 {
			DatosArch += "B" + fmt.Sprint(pos) + "->" + "I" + fmt.Sprint(int(pos2)-1) + "\n"
			DatosArch += TreeInodoArch((int64(pos2) - 1), dsk, SuperBlock)

		}
	}

	return DatosArch
}

func TreeInodoArch(pos int64, dsk *os.File, SuperBlock Estructuras.Sblock) string {

	InodoArchivo := ReadInode(Estructuras.I_node{}, SuperBlock.S_inode_start+SuperBlock.S_inode_size*pos, dsk)

	DatosArch := "I" + fmt.Sprint(pos) + "[label=<\n<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"1\" CELLPADDING=\"4\">\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"purple\" COLSPAN=\"2\">I" + fmt.Sprint(pos) + "</TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_uid</TD><TD>" + fmt.Sprint(InodoArchivo.I_uid) + " </TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_gid</TD><TD>" + fmt.Sprint(InodoArchivo.I_gid) + " </TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_size</TD><TD>" + fmt.Sprint(InodoArchivo.I_size) + " </TD>\n</TR>\n"
	DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">I_atime</TD><TD>" + string(InodoArchivo.I_atime[:]) + " </TD>\n</TR>\n"

	for i := 0; i < 16; i++ {
		act := InodoArchivo.I_block[i]
		if (int64(act) - 1) != -1 {
			DatosArch += "<TR>\n<TD BGCOLOR=\"lightblue\">Ap" + fmt.Sprint(i) + "</TD><TD>" + fmt.Sprint(act) + " </TD>\n</TR>\n"
		}
	}
	DatosArch += "</TABLE>>];\n"

	if string(InodoArchivo.I_type) == "1" {
		DatosArch += TreeFileBlocks(InodoArchivo, dsk, SuperBlock, pos)
	} else {
		for i := 0; i < 16; i++ {
			act := InodoArchivo.I_block[i]
			if (int64(act)-1) != -1 && (int64(act)-1) != 0 {
				DatosArch += "I" + fmt.Sprint(pos) + "->" + "B" + fmt.Sprint((int64(act) - 1)) + "\n"
				DatosArch += TreeBlckCarpetas((int64(act) - 1), dsk, SuperBlock)
			}
		}
	}

	return DatosArch

}

func TreeFileBlocks(Inodo Estructuras.I_node, dsk *os.File, SuperBlock Estructuras.Sblock, Father int64) string {
	DatosArch := ""

	for i := 0; i < 16; i++ {
		act := Inodo.I_block[i]
		if (int64(act) - 1) != -1 {
			is := (int64(act) - 1)
			DatosArch += "I" + fmt.Sprint(Father) + "->" + "B" + fmt.Sprint((int64(act) - 1)) + "\n"
			DatosArch += "B" + fmt.Sprint((int64(act) - 1)) + "[label=<\n<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"1\" CELLPADDING=\"4\">\n"
			DatosArch += "<TR>\n<TD BGCOLOR=\"purple\">Bloque" + fmt.Sprint((int64(act) - 1)) + "</TD>\n</TR>\n"
			pos := SuperBlock.S_block_start + (int64(unsafe.Sizeof(Estructuras.BloqueArchivos{})) * is)
			bl_archivo := ReadBloqueArchivo(Estructuras.BloqueArchivos{}, pos, dsk)
			DatosArch += "<TR>\n<TD>" + remove0s((bl_archivo.B_content[:])) + "</TD>\n</TR>\n"
			DatosArch += "</TABLE>>];\n"
		}
	}

	return DatosArch
}

func remove0s(Data []byte) string {
	Cleanse := ""
	for i := 0; i < len(Data); i++ {
		if Data[i] != 0 {
			Cleanse += string(Data[i])
		}
	}

	return Cleanse
}
