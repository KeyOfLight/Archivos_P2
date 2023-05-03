package operaciones

import (
	"Proyecto2/Estructuras"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
)

func Reportes(parameters Estructuras.ParamStruct) {

	if strings.Contains(parameters.Nombre, "disk") {
		RepDisk(parameters)
	} else if strings.Contains(parameters.Nombre, "tree") {
		RepTree(parameters)
	} else if strings.Contains(parameters.Nombre, "file") {
		RepFile(parameters)
	} else if strings.Contains(parameters.Nombre, "sb") {
		RepSb(parameters)
	}
}

func RepTree(parameters Estructuras.ParamStruct) {

	var ListaMontados = Estructuras.ListaMontados{}
	ListaMontados = (&ListaMontados).GetLista()
	Uss := SearchID(parameters.Id, ListaMontados)

	StartPoint := Uss.StartPoint
	path := Uss.Path

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	mbr := OpenMBR(path)
	print(mbr.Dsk_fit)
	Diagrama := "digraph Tree{\n node [shape=plaintext];\nrankdir=LR;\n "

	SuperBlock := ReadSBlock(Estructuras.Sblock{}, StartPoint, dsk)
	Diagrama += MkTree(SuperBlock, dsk)

	Diagrama += "}"

	crearDirectorio(parameters.Direccion)
	d1 := []byte(Diagrama)
	err3 := os.WriteFile(parameters.Direccion, d1, 0644)
	if err3 != nil {
		panic(err3)
	}

	Limpio := RemoveDot(parameters.Direccion)

	cmd := exec.Command("dot", parameters.Direccion, "-Tpdf", "-o", Limpio+".pdf")
	errcmd := cmd.Run()
	if errcmd != nil {
		fmt.Println("error", errcmd.Error())
		return
	}

}

func RepFile(parameters Estructuras.ParamStruct) {

	var ListaMontados = Estructuras.ListaMontados{}
	ListaMontados = (&ListaMontados).GetLista()
	Uss := SearchID(parameters.Id, ListaMontados)
	StartPoint := Uss.StartPoint
	path := Uss.Path

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	Diagrama := "digraph File{\n node [shape=plaintext];\nrankdir=LR;\n "
	Diagrama += "SBlock[label=<\n<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"1\" CELLPADDING=\"4\">\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"purple\">File</TD>\n</TR>\n"

	SuperBlock := ReadSBlock(Estructuras.Sblock{}, StartPoint, dsk)
	Info := LeerArchivoMkfs(SuperBlock, parameters.Pwd, dsk)

	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">" + string(Info) + " </TD>\n</TR>\n"
	Diagrama += "</TABLE>>];\n"

	Diagrama += "}"
	crearDirectorio(parameters.Direccion)
	d1 := []byte(Diagrama)
	err3 := os.WriteFile(parameters.Direccion, d1, 0644)
	if err3 != nil {
		panic(err3)
	}

	Limpio := RemoveDot(parameters.Direccion)

	cmd := exec.Command("dot", parameters.Direccion, "-Tpdf", "-o", Limpio+".pdf")
	errcmd := cmd.Run()
	if errcmd != nil {
		fmt.Println("error", errcmd.Error())
		return
	}
}

func RepSb(parameters Estructuras.ParamStruct) {
	var ListaMontados = Estructuras.ListaMontados{}
	ListaMontados = (&ListaMontados).GetLista()
	Uss := SearchID(parameters.Id, ListaMontados)
	path := Uss.Path
	StartPoint := Uss.StartPoint

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	SuperBlock := ReadSBlock(Estructuras.Sblock{}, StartPoint, dsk)
	Diagrama := "digraph SBloques{\n node [shape=plaintext];\nrankdir=LR;\n "
	Diagrama += "SBlock[label=<\n<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"1\" CELLPADDING=\"4\">\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"purple\" COLSPAN=\"2\">Super Bloque</TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_filesystem_type</TD><TD>" + fmt.Sprint(SuperBlock.S_filesystem_type) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_inodes_count</TD><TD>" + fmt.Sprint(SuperBlock.S_inodes_count) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_blocks_count</TD><TD>" + fmt.Sprint(SuperBlock.S_blocks_count) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_free_blocks_count</TD><TD>" + fmt.Sprint(SuperBlock.S_free_blocks_count) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_free_inodes_count</TD><TD>" + fmt.Sprint(SuperBlock.S_free_inodes_count) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_mtime</TD><TD>" + string(SuperBlock.S_mtime[:]) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_mnt_count</TD><TD>" + fmt.Sprint(SuperBlock.S_mnt_count) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_magic</TD><TD>" + fmt.Sprint(SuperBlock.S_magic) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_inode_size</TD><TD>" + fmt.Sprint(SuperBlock.S_inode_size) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_block_size</TD><TD>" + fmt.Sprint(SuperBlock.S_block_size) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_firts_ino</TD><TD>" + fmt.Sprint(SuperBlock.S_firts_ino) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_first_blo</TD><TD>" + fmt.Sprint(SuperBlock.S_first_blo) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_bm_inode_start</TD><TD>" + fmt.Sprint(SuperBlock.S_bm_inode_start) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_bm_block_start</TD><TD>" + fmt.Sprint(SuperBlock.S_bm_block_start) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_inode_start</TD><TD>" + fmt.Sprint(SuperBlock.S_inode_start) + " </TD>\n</TR>\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"lightblue\">s_block_start</TD><TD>" + fmt.Sprint(SuperBlock.S_block_start) + " </TD>\n</TR>\n"
	Diagrama += "</TABLE>>];\n"

	Diagrama += "}"

	crearDirectorio(parameters.Direccion)
	dsk, err2 := os.Create(parameters.Direccion)
	defer dsk.Close()
	if err2 != nil {
		fmt.Println("No se pudo crear el archivo deseado")
		fmt.Println(err2.Error())
		return
	}

	d1 := []byte(Diagrama)
	err3 := os.WriteFile(parameters.Direccion, d1, 0644)
	if err3 != nil {
		panic(err3)
	}

	Limpio := RemoveDot(parameters.Direccion)

	cmd := exec.Command("dot", parameters.Direccion, "-Tpdf", "-o", Limpio+".pdf")
	errcmd := cmd.Run()
	if errcmd != nil {
		fmt.Println("error", errcmd.Error())
		return
	}
}

func RepDisk(parammeters Estructuras.ParamStruct) {
	var ListaMontados = Estructuras.ListaMontados{}
	ListaMontados = (&ListaMontados).GetLista()
	Uss := SearchID(parammeters.Id, ListaMontados)
	path := Uss.Path

	dsk, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer dsk.Close()

	if err != nil {
		fmt.Println("No se pudo encontrar el archivo deseado")
		os.Exit(1)
	}

	Diagrama := "digraph DSK{\n node [shape=plaintext];\n "

	MBR := OpenMBR(path)

	Diagrama += "struct3 [label=<\n<TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"1\" CELLPADDING=\"4\">\n"
	Diagrama += "<TR>\n<TD BGCOLOR=\"purple\" ROWSPAN=\"2\">MBR</TD>\n"

	ocupado := int64(0)
	ExtPart := Estructuras.Particion{}
	for i := 0; i < 4; i++ {
		Act := MBR.Mbr_partition[i]
		if Act.Part_status == '1' {
			if Act.Part_type == 'e' {
				ExtPart = MBR.Mbr_partition[i]
				ocupado += int64(Act.Part_size)
				Externos := ContarEBRs(Act.Part_start, dsk, 0)
				Diagrama += "<TD COLSPAN=\"" + fmt.Sprint(Externos*3) + "\">Extendida</TD>\n"
			} else {
				var porcentaje float64
				porcentaje = float64(Act.Part_size) / float64(MBR.Mbr_tamano)
				porcentaje = porcentaje * 100
				ocupado += (Act.Part_size)
				Diagrama += "<TD ROWSPAN=\"2\">Primaria <BR/>" + fmt.Sprint(math.Round(porcentaje)) + "%</TD>\n"
			}

		} else {
			procentaje := (float64(MBR.Mbr_tamano) - float64(ocupado)) / float64(MBR.Mbr_tamano)
			procentaje = procentaje * 100
			Diagrama += "<TD ROWSPAN=\"2\">Libre <BR/>" + fmt.Sprint(math.Round(procentaje)) + "%</TD>\n"
		}
	}

	Diagrama += "</TR>\n"
	Diagrama += "<TR>\n"
	Diagrama += DiskEBR(ExtPart.Part_start, MBR, dsk, ExtPart.Part_size, 0)
	Diagrama += "</TR>\n</TABLE>>];\n}"

	crearDirectorio(parammeters.Direccion)
	f, err := os.Create(parammeters.Direccion)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.WriteString(Diagrama + "\n")
	if err != nil {
		log.Fatal(err)
	}

	Limpio := RemoveDot(parammeters.Direccion)

	cmd := exec.Command("dot", parammeters.Direccion, "-Tpdf", "-o", Limpio+".pdf")
	errcmd := cmd.Run()
	if errcmd != nil {
		fmt.Println("error", errcmd.Error())
		return
	}

	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()

}

func DiskEBR(Part int64, mbr Estructuras.MBR, dsk *os.File, total int64, Written int64) string {
	StartPoint := Part
	tempebr := ReadEbr(StartPoint, dsk)
	escrito := Written
	Diagrama := ""

	if tempebr.Part_status == '1' {
		escrito += tempebr.Part_size
		porcentaje := float64(tempebr.Part_size) / float64(mbr.Mbr_tamano)
		porcentaje = porcentaje * 100
		Diagrama += "<TD BGCOLOR=\"lightblue\" >EBR</TD><TD>Logica <BR/> " + fmt.Sprint(math.Round(porcentaje)) + " % </TD>"
		Diagrama += DiskEBR(tempebr.Part_next, mbr, dsk, total, escrito)
	} else {
		porcentaje := float64(total - escrito)
		porcentaje = porcentaje / float64(mbr.Mbr_tamano)
		porcentaje = porcentaje * 100
		Diagrama += "<TD BGCOLOR=\"lightblue\" >EBR</TD><TD>Libre  <BR/> " + fmt.Sprint(math.Round(porcentaje)) + " % </TD>"
	}

	return Diagrama
}

func ContarEBRs(StartPoint int64, dsk *os.File, Ocupado int64) int64 {
	tempebr := ReadEbr(StartPoint, dsk)

	if tempebr.Part_status == '1' {
		Ocupado += 1
		return ContarEBRs(tempebr.Part_next, dsk, Ocupado)
	}
	return Ocupado
}

func RemoveDot(Ruta string) string {
	retorno := ""
	SeparadoLineas := strings.Split(Ruta, ".dot")
	for _, linea := range SeparadoLineas {
		retorno += linea
	}

	return retorno
}

func SearchID(ID string, Part Estructuras.ListaMontados) Estructuras.PartMounted {
	Empt := Estructuras.PartMounted{}
	for i := 0; i < len(Part.Montado); i++ {
		if Part.Montado[i].Id == ID {
			return Part.Montado[i]
		}
	}

	return Empt
}
