package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
	"tailscale.com/tsnet"
)

var (
	addr     = flag.String("addr", ":80", "address to listen on")
	hostname = flag.String("hostname", "wol", "hostname to use on the tailnet")
)

func serve() {

	k := os.Getenv("TS_AUTHKEY")
	if k == "" {
		appPath, err := appPath()
		if err != nil {
			log.Fatal(err)
		}

		dotenvPath := filepath.Join(appPath, ".env")
		err = godotenv.Load(dotenvPath)
		if err != nil {
			slog.Info("Cannot load .env file")
		}
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

		c, err := getConfig()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		hosts := c.Hosts

		a := r.URL.Query().Get("a")
		if a == "" {
			http.Error(w, "mac address is required", 400)
		}
		a = strings.Replace(a, "-", ":", -1)

		t := findTarget(a, c)
		if t == nil {
			t = &ConfigTarget{Mac: a}
		}

		for _, host := range hosts {

			var cmd string

			if t.PreferredCommand == nil {
				cmd = buildNetCatCommand(a, t.Ip, t.Port)
			} else if *t.PreferredCommand == "wol" {
				cmd = buildWakeOnLanCommand(a, t.Ip, t.Port)
			} else if *t.PreferredCommand == "netcat" {
				cmd = buildNetCatCommand(a, t.Ip, t.Port)
			} else {
				http.Error(w, "Invalid preferredCommand", 400)
				return
			}

			_, err := sendCommand(host.Host, host.User, cmd, host.Port, host.Password, host.Identity)
			if err == nil {
				// http.ResponseWriter.WriteHeader(w, 204)
				fmt.Printf("Sent WOL to %s\n", a)
				fmt.Printf("Command: %s", cmd)
				return
			}
		}

		http.Error(w, "Failed to send WOL", 500)
	})

	appPath, err := appPath()
	if err == nil {
		publicDir := filepath.Join(appPath, "public/")
		http.Handle("/", http.FileServer(http.Dir(publicDir)))
	}

	log.Fatal(http.Serve(ln, nil))
}

func firstLabel(s string) string {
	s, _, _ = strings.Cut(s, ".")
	return s
}

func sendCommand(host string, user string, command string, port *int, password *string, identityFile *string) (string, error) {

	c := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
	}

	if password != nil {
		c.Auth = append(c.Auth, ssh.Password(*password))
	}

	var p int
	if port != nil {
		p = *port
	} else {
		p = 22 // Default port value
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, p), c)
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

func buildNetCatCommand(mac string, broadcast_ip *string, port *int) string {
	var ip string
	if broadcast_ip != nil {
		ip = *broadcast_ip
	} else {
		ip = "255.255.255.255" // Default broadcast ip
	}

	var p int
	if port != nil {
		p = *port
	} else {
		p = 9 // Default port value
	}

	a := regexp.MustCompile("[:-]").ReplaceAllString(mac, "")

	tpl := "bash -c '(printf 'FF%%.0s' {1..6}; printf %s'%%.0s' {1..16}) | sed 's/../\\\\\\\\\\\\\\\\x&/g' | xargs printf '%%b' | netcat -u -b -w1 %s %d'"
	return fmt.Sprintf(tpl, a, ip, p)
}

func buildWakeOnLanCommand(mac string, broadcast_ip *string, port *int) string {
	var ip string
	if broadcast_ip != nil {
		ip = *broadcast_ip
	} else {
		ip = "255.255.255.255" // Default broadcast ip
	}

	var p int
	if port != nil {
		p = *port
	} else {
		p = 9 // Default port value
	}

	tpl := "wakeonlan -i %s -p %d %s"
	return fmt.Sprintf(tpl, ip, p, mac)
}
