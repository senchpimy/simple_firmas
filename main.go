package main
import (
        "log"
	"bufio"
	"fmt"
        "net/http"
        "os"
        "io/ioutil"
        "strings"
	"html/template"
)

const maxSize = 1 * 1024 * 1024; //one megaByte
const path = "firmas";
func main() {
        //LimitFile("firmas","hola",49999)

        fs := http.FileServer(http.Dir("."))
        http.Handle("/files", fs)

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				switch r.Method{
					case "GET":
						tpl, err := template.ParseFiles("index.html")
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
	
						data,err:=readLines()
						if err != nil {return}
	
						err = tpl.Execute(w, data)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					case "POST":
						fmt.Println("Mensaje Recibido")
						message:=r.FormValue("message")
						LimitFile(message)
							http.Redirect(w, r, "/", http.StatusSeeOther)
				}
			     })

        log.Print("Listening on :3001..")
        err := http.ListenAndServe(":3001", nil)
        if err != nil {
                log.Fatal(err)
        }
}

func LimitFile(newdata string) error {
	file:=path
        // get file size
        info, err := os.Stat(file)
        if err != nil {
                return err
        }
        size := info.Size()
        if size > maxSize {
                // read the file
                data, err := ioutil.ReadFile(file)
                if err != nil {
                        return err
                }
                // split the file into lines
                lines := strings.Split(string(data), "\n")
                // delete the first line
                lines = lines[1:]
                // write the remaining lines back to the file
                lines = append(lines, newdata+"\n")
                err = ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
                if err != nil {
                        return err
                }
        }else{
                file, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
                if err != nil {
                        fmt.Println(err)
                        return err
                }
                defer file.Close()

                // write to the file
                _, err = file.WriteString(newdata+"\n")
                if err != nil {
                        fmt.Println(err)
                        return err
                }
        }

        return nil
}

func readLines() ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
