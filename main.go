package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"bytes"
	_ "errors"
)

func stringTab(nbr int) (str string) {
	for i := 0; i < nbr; i++ {
		str += "\t"
	}
	return str
}

func includeContent(buf []byte, f *os.File, tab *int) (er error) {
	sbuf := bytes.Split(buf, []byte("\",")) // Split with ", to simple key with value
	lbuf := len(sbuf)
	for i := 0; i < lbuf; i++ {
		if len(sbuf[i]) > 0 {
			_, er = f.WriteString(stringTab(*tab))
			if er != nil {
				return er
			}
			_, er = f.Write(sbuf[i])
			if er != nil {
				return er
			}
			_, er = f.WriteString("\n")
			if er != nil {
				return er
			}
		}
	}
	return er
}

func closeBracket(tab *int, buf []byte, f *os.File) (er error) {
	sbuf := bytes.Split(buf, []byte("}"))
	lbuf := len(sbuf)
	for i := 0; i < lbuf; i++ {
		er = includeContent(sbuf[i], f, tab)
		if er != nil {
			return er
		}
		if lbuf > 1 && i + 1 < lbuf {
			*tab--
			_, er = f.WriteString(stringTab(*tab))
			if er != nil {
				return er
			}
			_, er = f.WriteString("}\n")
			if er != nil {
				return er
			}
		}
	}
	return er
}

func constructor(buf []byte, f *os.File) (er error) {
	tab := 0
	sbuf := bytes.Split(buf, []byte(string('{')))
	lbuf := len(sbuf)
	for i := 0; i < lbuf; i++ {
		_, er = f.WriteString(stringTab(tab))
		if er != nil {
			return er
		}
		_, er = f.WriteString("{\n")
		if er != nil {
			return er
		}
		tab++
		er = closeBracket(&tab, sbuf[i], f)
		if er != nil {
			return er
		}
	}
	return er
}

func main () {
	if len(os.Args) != 3 {
		fmt.Println("./Usage 'path/namefile.json' 'path/newfileShema'")
		return
	}
	buf, er := ioutil.ReadFile(os.Args[1])
	if er != nil {
		fmt.Println("Can't read file specified: ", er)
		return
	}
	F, er := os.OpenFile(os.Args[2], os.O_CREATE | os.O_TRUNC | os.O_RDWR, 0666)
	if er != nil {
		fmt.Println("Can't create a new file: ", er)
		return
	}
	defer F.Close()
	er = constructor(buf, F)
	if er != nil {
		fmt.Println("Error on Parsing file: ", er)
		return
	}
	return
}