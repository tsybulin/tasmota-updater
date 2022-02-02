package tasmotaupdater

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
var HTTP_TIMEOUT_SEC = 30 * time.Second

func ping(url string, what string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 200 && strings.Contains(body.String(), what), nil
}

func createFormFile(w *multipart.Writer, fieldname, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, quoteEscaper.Replace(fieldname), quoteEscaper.Replace(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func sendFile(url string, file string, what string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer f.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := createFormFile(w, "u2", file, "application/x-gzip")
	if err != nil {
		return false, err
	}
	_, err = io.Copy(fw, f)
	if err != nil {
		return false, err
	}

	err = w.Close()
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return false, err
	}

	err = resp.Body.Close()
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 200 && strings.Contains(body.String(), what), nil
}

func upload(ip string, file string) error {
	log.Print("Ping Main")
	ok, err := ping("http://"+ip+"/", "Firmware Upgrade")
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("main page not found")
	}

	log.Print("Ping Up")
	ok, err = ping("http://"+ip+"/up", "Upgrade by file upload")
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("update page not found")
	}

	log.Print("Upload ", file)
	url := "http://" + ip + "/u2"
	ok, err = sendFile(url, file, "Successful")
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("upload error")
	}

	return nil
}

func update(tasmota tasmota) {
	log.Print("Update ", tasmota.Ip, " ", tasmota.Name, " from ", tasmota.Version)

	err := upload(tasmota.Ip, "tasmota-minimal.bin.gz")
	if err != nil {
		log.Print(err)
		return
	}

	time.Sleep(HTTP_TIMEOUT_SEC)

	err = upload(tasmota.Ip, "tasmota.bin.gz")
	if err != nil {
		log.Print(err)
		return
	}

	time.Sleep(HTTP_TIMEOUT_SEC)

	log.Print("Check Main")
	ok, err := ping("http://"+tasmota.Ip+"/", "Firmware Upgrade")
	if err != nil {
		log.Print(err)
		return
	}

	if !ok {
		log.Print("main page not found")
		return
	}

	log.Print("--- done ---")
}

func Update(version string) {
	for ip, discovery := range tasmotas {
		if discovery.Version == version {
			log.Print("Skip ", ip, " ", discovery.Name, " on ", discovery.Version)
		} else {
			update(discovery)
		}
	}
}
