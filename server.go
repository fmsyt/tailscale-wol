package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
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

func serve() {

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

		appConfig, err := getConfig()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		hosts := appConfig.Hosts

		mac_address := r.URL.Query().Get("a")
		if mac_address == "" {
			http.Error(w, "mac address is required", 400)
		}
		mac_address = strings.Replace(mac_address, "-", ":", -1)

		command := fmt.Sprintf("/usr/bin/wakeonlan %s", mac_address)

		if len(hosts) == 0 {
			http.Error(w, "host is not defined", 400)
			return
		}

		host := hosts[0]

		_, err = wol(host.Host, host.User, command, host.Port, host.Password, host.Identity)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	})

	http.Handle("/", http.FileServer(http.Dir("public/")))

	log.Fatal(http.Serve(ln, nil))
}

func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}

func wol(host string, user string, command string, port *int, password *string, identityFile *string) (string, error) {

	config := &ssh.ClientConfig{
		User: user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
	}

	if password != nil {
		config.Auth = append(config.Auth, ssh.Password(*password))
	}

	p := *port
	if port == nil {
		p = 22
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, p), config)
	if err != nil {
		return "", err
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(command); err != nil {
		return "", err
	}

	return b.String(), nil
}
