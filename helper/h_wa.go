package helper

import (
	"encoding/gob"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	qrterm "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Jeffail/gabs"
	whatsapp "github.com/dimaskiddo/go-whatsapp"
)

type WAHandler struct {
	WAConn        *whatsapp.Conn
	JID           string
	StartTime     uint64
	ReconnectTime int
	TestMode      bool
}

func (wah *WAHandler) HandleError(err error) {
	_, eMatch := err.(*whatsapp.ErrConnectionFailed)
	if eMatch {
		log.Println("error connection: connection closed will reconnecting after " + strconv.Itoa(wah.ReconnectTime) + " seconds")

		wah.StartTime = uint64(time.Now().Unix())
		<-time.After(time.Duration(wah.ReconnectTime) * time.Second)

		err := wah.WAConn.Restore()
		if err != nil {
			log.Println(strings.ToLower(err.Error()))
		}
	} else {
		if strings.Contains(strings.ToLower(err.Error()), "server closed connection") {
			return
		}

		log.Println(strings.ToLower(err.Error()))
	}
}

func (wah *WAHandler) HandleTextMessage(msg whatsapp.TextMessage) {
	if !strings.Contains(strings.ToLower(msg.Text), "@bot") || msg.Info.Timestamp < wah.StartTime {
		return
	}

	msgRecieved := strings.SplitN(strings.ToLower(msg.Text), " ", 2)[1]
	log.Printf("recieved text message\nTimestamp:\t%v\nMessage ID:\t%v\nQuoted to ID:\t%v\nRemote JID:\t%v\nMessage:\t%v\n", msg.Info.Timestamp, msg.Info.Id, msg.Info.QuotedMessageID, msg.Info.RemoteJid, msgRecieved)

	if wah.TestMode {
		if msg.Info.FromMe && msg.Info.RemoteJid == wah.JID {
			msgResponse := WACmd.Path("data." + msgRecieved + ".res").Data()

			err := WAMessageText(wah.WAConn, msg.Info.RemoteJid, msgResponse.(string), 0)
			if err != nil {
				log.Println("error while sending message: " + err.Error())
			}
		}
	} else {
		msgResponse := WACmd.Path("data." + msgRecieved + ".res").Data()

		err := WAMessageText(wah.WAConn, msg.Info.RemoteJid, msgResponse.(string), 0)
		if err != nil {
			log.Println("error while sending message: " + err.Error())
		}
	}
}

var WACmd *gabs.Container
var WAConn *whatsapp.Conn

func WAInitCmd(file string) (*gabs.Container, error) {
	json, err := gabs.ParseJSONFile(file)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func WAInitConn(timeout int) (*whatsapp.Conn, error) {
	conn, err := whatsapp.NewConn(time.Duration(timeout) * time.Second)
	if err != nil {
		return nil, err
	}
	conn.SetClientName("Go WhatsApp CLI", "Go WhatsApp")
	return conn, nil
}

func WASessionLoad(file string) (whatsapp.Session, error) {
	session := whatsapp.Session{}

	buffer, err := os.Open(file)
	if err != nil {
		return session, err
	}
	defer buffer.Close()

	err = gob.NewDecoder(buffer).Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func WASessionSave(file string, session whatsapp.Session) error {
	buffer, err := os.Create(file)
	if err != nil {
		return err
	}
	defer buffer.Close()

	err = gob.NewEncoder(buffer).Encode(session)
	if err != nil {
		return err
	}

	return nil
}

func WASessionExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	return true
}

func WASessionLogin(conn *whatsapp.Conn, file string) error {
	if conn != nil {
		defer func() {
			_, _ = conn.Disconnect()
		}()

		if WASessionExist(file) {
			return errors.New("error session: session file already exist, please logout first")
		}

		qrstr := make(chan string)
		go func() {
			term := qrterm.New()
			term.Get(<-qrstr).Print()
		}()

		session, err := conn.Login(qrstr)
		if err != nil {
			return err
		}

		err = WASessionSave(file, session)
		if err != nil {
			return err
		}
	} else {
		return errors.New("error connection: connection is not valid")
	}

	return nil
}

func WASessionRestore(conn *whatsapp.Conn, file string) error {
	if conn != nil {
		if !WASessionExist(file) {
			_, _ = conn.Disconnect()
			return errors.New("error session: session file doesn't exist, please login first")
		}

		session, err := WASessionLoad(file)
		if err != nil {
			_ = os.Remove(file)
			_, _ = conn.Disconnect()
			return errors.New("error session: session not valid, removing session file")
		}

		session, err = conn.RestoreWithSession(session)
		if err != nil {
			_, _ = conn.Disconnect()
			return err
		}
	} else {
		return errors.New("error connection: connection is not valid")
	}

	return nil
}

func WASessionLogout(conn *whatsapp.Conn, file string) error {
	if conn != nil {
		defer func() {
			_, _ = conn.Disconnect()
		}()

		err := conn.Logout()
		if err != nil {
			return err
		}

		_ = os.Remove(file)
	} else {
		return errors.New("error connection: connection is not valid")
	}

	return nil
}

func WAMessageText(conn *whatsapp.Conn, msgJID string, msgText string, msgDelay int) error {
	if conn != nil {
		content := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: msgJID,
			},
			Text: msgText,
		}

		<-time.After(time.Duration(msgDelay) * time.Second)

		err := conn.Send(content)
		if err != nil {
			return err
		}
	} else {
		return errors.New("error connection: connection is not valid")
	}

	return nil
}
