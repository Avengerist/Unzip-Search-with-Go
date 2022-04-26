package main
import (
	"fmt"
	"path/filepath"
	"io/fs"
	"io"
	"os"
	"archive/zip"
	"encoding/csv"
)
func Search (path string) error{
	file,err := os.Open(path)
	var slice [][]string
	if err != nil{
		return err
	}
	defer file.Close()
	r:= csv.NewReader(file)
	for {
		row, err:= r.Read()
		if err == io.EOF{
			break
		}
		if err != nil{
			return err
		}
		slice=append(slice,row)
	}
	if (len(slice[0])==10){
		fmt.Printf("File: %s \nValue of Row 5 and Cell 3 is: %s\n",path,slice[4][2])
		os.Exit(0)
	}
	return nil

}
func Unzip(f *zip.File, path string) error{
	filePath:=filepath.Join(f.Name)
	if f.FileInfo().IsDir(){
		if err:=os.MkdirAll(filePath,os.ModePerm); err != nil{
			fmt.Println(err)
			return err
		}
		if err:=os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil{
			return err
		}
	}
	fmt.Println(filePath)
	destinationFile,err := os.OpenFile(filePath,os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil{
		return  err
	}
	zippedFile,err := f.Open()
	if err != nil{
		return err
	}
	defer zippedFile.Close()
	if _,err := io.Copy(destinationFile,zippedFile); err != nil{
		return err
	}
	if f.FileInfo().IsDir() == false{
		Search(filePath)
	}
	defer destinationFile.Close()
	return nil
}
func GetContent(path string) error{
	r,err := zip.OpenReader(path)
	defer r.Close()
	if err != nil{
		fmt.Println(err)
		return err
	}
	for _,f := range r.File{
		Unzip(f, path)
	}
	return nil

}
func WalkFunc(path string, info fs.FileInfo, err error) error{
	fmt.Println(info.Name(),":",filepath.Ext(path),":",path)
	if err != nil{
		fmt.Println(err)
		return err
	}
	if filepath.Ext(path) == ".zip" {
		GetContent(path)
	}
	return nil


}
func main(){
	filepath.Walk(".",WalkFunc)

}
