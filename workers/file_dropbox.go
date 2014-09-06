package workers

import (
	"fmt"
	"github.com/callumj/auto-convert/shared"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func processFile(fReq FileRequest) {
	extension := filepath.Ext(fReq.Path)
	extension = strings.TrimPrefix(extension, ".")
	matchingCmd, found := shared.Config.Extensions[extension]
	if !found {
		log.Printf("%v is not configured as a extension. (%v)\n", extension, fReq.Path)
		return
	}

	url := fmt.Sprintf("https://api-content.dropbox.com/1/files/auto/%s", fReq.Path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	acc := shared.FetchAccount(shared.Account{Uid: fReq.Uid})
	if acc == nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", acc.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	} else {
		defer resp.Body.Close()
		tempFile, err := ioutil.TempFile("", extension)
		if err != nil {
			log.Println(err)
			return
		}
		defer os.Remove(tempFile.Name())

		io.Copy(tempFile, resp.Body)

		outTempFile, err := ioutil.TempFile("", matchingCmd.Extension)
		if err != nil {
			log.Println(err)
			return
		}
		defer os.Remove(outTempFile.Name())
		confedCmd := strings.Replace(matchingCmd.Cmd, "{{file}}", tempFile.Name(), 1)
		confedCmd = strings.Replace(confedCmd, "{{out_file}}", outTempFile.Name(), 1)

		log.Println(confedCmd)
		cmd := exec.Command("sh", "-c", confedCmd)
		err = cmd.Start()
		if err != nil {
			log.Println(err)
		}
		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		if err != nil {
			log.Println(err)
			return
		} else {
			newPath := strings.Replace(fReq.Path, fReq.MatchedPath, shared.Config.Maps[fReq.MatchedPath], 1)
			newPath = strings.TrimSuffix(newPath, fmt.Sprintf(".%s", extension))
			newPath = fmt.Sprintf("%s.%s", newPath, matchingCmd.Extension)
			url := fmt.Sprintf("https://api-content.dropbox.com/1/files_put/auto/%s", newPath)

			uploadReq, err := http.NewRequest("PUT", url, outTempFile)
			uploadReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", acc.Token))

			if err != nil {
				log.Println(err)
				return
			} else {
				resp, err = client.Do(uploadReq)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(resp.Status)
				}
			}
		}
	}
}
