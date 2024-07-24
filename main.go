package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log/slog"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/ssh"
	"tailscale.com/tsnet"
	"github.com/joho/godotenv"
)

var (
	addr     = flag.String("addr", ":80", "address to listen on")
	hostname = flag.String("hostname", "wol", "hostname to use on the tailnet")
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Info("Cannot load .env file")
	}

	flag.Parse()
	s := new(tsnet.Server)
	s.Hostname = *hostname
	defer s.Close()
	ln, err := s.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/whoami/", func(w http.ResponseWriter, r *http.Request) {
		who, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprintf(w, "<html><body><h1>Hello, tailnet!</h1>\n")
		fmt.Fprintf(w, "<p>You are <b>%s</b> from <b>%s</b> (%s)</p>",
		html.EscapeString(who.UserProfile.LoginName),
		html.EscapeString(firstLabel(who.Node.ComputedName)),
		r.RemoteAddr)
	})

	http.HandleFunc("/run/", func(w http.ResponseWriter, r *http.Request) {
		_, err := lc.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// TODO: 認証設定
		config := &ssh.ClientConfig{
			User: "",
			Auth: []ssh.AuthMethod{
				ssh.Password(""),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
		}

		client, err := ssh.Dial("tcp", "rpi:22", config)
		if err != nil {
			log.Fatal("Failed to dial: ", err)
		}

		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			log.Fatal("Failed to create session: ", err)
		}

		defer session.Close()

		var b bytes.Buffer
		session.Stdout = &b

		mac_address := r.URL.Query().Get("a")
		if mac_address == "" {
			http.Error(w, "mac address is required", 400)
		}
		mac_address = strings.Replace(mac_address, "-", ":", -1)

		command := fmt.Sprintf("/usr/bin/wakeonlan %s", mac_address)

		if err := session.Run(command); err != nil {
			log.Fatal("Failed to run: " + err.Error())
		}

		fmt.Fprint(w, b.String())
	})

	http.Handle("/", http.FileServer(http.Dir("public/")))

	log.Fatal(http.Serve(ln, nil))
}

func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}

type Config struct {
	Hosts []ConfigHost
}

type ConfigHost struct {
	host     string
	user     string
	port     *int
	password *string
	identity *string
}

func readConfigFromFile() (*Config, error) {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)

	return &config, err
}
