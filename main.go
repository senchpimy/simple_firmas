package main
import (
        "log"
        "net/http"
        "os"
        "io/ioutil"
        "strings"
)

func main() {
        LimitFile("firmas","hola",49999)
        fs := http.FileServer(http.Dir("."))
        http.Handle("/", fs)

        log.Print("Listening on :3001..")
        err := http.ListenAndServe(":3001", nil)
        if err != nil {
                log.Fatal(err)
        }
}

func LimitFile(file, newdata string, maxSize int64) error {
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
                lines = appens(lines, newdata+"\n")
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

                }

        return nil
}

