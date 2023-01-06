package utilities

import (
	"bufio"
	"bytes"
	"exportcsv/models"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	sshConn  = map[string]*ssh.Client{}
	mu       sync.Mutex
	Creation = false
	Panorama = models.Device{}
)

func PolicyCreation() models.Device {

	fmt.Print("Do you want to continue building policies in Panorama? Y or N :")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := scanner.Text()

	if (response == "Y" || response == "y") && !Creation {

		fmt.Print("Enter Panorama Hostname/IP Address: ")

		scanner.Scan()
		hostname := scanner.Text()

		fmt.Print("Enter Panorama Username: ")

		scanner.Scan()
		username := scanner.Text()

		fmt.Print("Enter Panorama Password: ")

		scanner.Scan()
		password := scanner.Text()

		fmt.Print("Enter DeviceGroup Name: ")

		scanner.Scan()
		devicegroup := scanner.Text()

		fmt.Print("Enter Template Name: ")

		scanner.Scan()
		template := scanner.Text()

		Creation = true

		Panorama = models.Device{
			Hostname:     hostname,
			Username:     username,
			Password:     password,
			DeviceGroup:  devicegroup,
			TemplateName: template,
		}
		return Panorama
	} else {
		return models.Device{}
	}
}

func ConnectSSH(user string, pass string, host string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conn, err
}

//WriteConn writes a command over ssh connection
func WriteConn(deviceIP string, username, pass, cmd string) (string, error) {
	if _, ok := sshConn[deviceIP]; ok {
		sess, err := sshConn[deviceIP].NewSession()
		if err != nil {
			log.Printf("%s Trying to connect again", err)
			mu.Lock()
			sshConn[deviceIP], err = ConnectSSH(username, pass, deviceIP)
			mu.Unlock()
			if err != nil {
				mu.Lock()
				delete(sshConn, deviceIP)
				mu.Unlock()
				log.Println(err)
				return "0", err
			}
			sess, err = sshConn[deviceIP].NewSession()
		}
		if err != nil {
			log.Println(err)
			return "0", err
		}
		defer sess.Close()

		stdin, err := sess.StdinPipe()
		if err != nil {
			log.Println(err)
			return "0", err
		}

		var b bytes.Buffer
		sess.Stdout = &b

		err = sess.Shell()
		if err != nil {
			log.Println(err)
			return "0", err
		}
		_, err = fmt.Fprintf(stdin, "%s\nexit\n", cmd)
		if err != nil {
			log.Println(err)
			return "0", err
		}

		err = sess.Wait()
		if err != nil {
			log.Println(err)
			return "0", err
		}

		space := regexp.MustCompile(`\s+`)
		str := space.ReplaceAllString(b.String(), " ")
		return str, nil
	}

	var err error
	mu.Lock()
	sshConn[deviceIP], err = ConnectSSH(username, pass, deviceIP)
	mu.Unlock()
	if err != nil {
		mu.Lock()
		delete(sshConn, deviceIP)
		mu.Unlock()
		log.Println(err)
		return "0", err
	}
	sess, err := sshConn[deviceIP].NewSession()
	if err != nil {
		return "0", err
	}
	defer sess.Close()

	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Println(err)
		return "0", err
	}

	var b bytes.Buffer
	sess.Stdout = &b

	err = sess.Shell()
	if err != nil {
		log.Println(err)
		return "0", err
	}
	_, err = fmt.Fprintf(stdin, "%s\nexit\n", cmd)
	if err != nil {
		log.Println(err)
		return "0", err
	}

	err = sess.Wait()
	if err != nil {
		log.Println(err)
		return "0", err
	}

	space := regexp.MustCompile(`\s+`)
	str := space.ReplaceAllString(b.String(), " ")
	return str, nil
}
