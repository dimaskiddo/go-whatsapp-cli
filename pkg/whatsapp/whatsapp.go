package whatsapp

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	qrterm "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	waproto "github.com/Rhymen/go-whatsapp/binary/proto"

	"github.com/dimaskiddo/go-whatsapp-cli/pkg/log"
	"github.com/dimaskiddo/go-whatsapp-cli/pkg/parser"
)

var WAConn *whatsapp.Conn
var WACMutex *sync.Mutex

type WAHandler struct {
	SessionConn   *whatsapp.Conn
	SessionJid    string
	SessionTag    string
	SessionFile   string
	SessionStart  uint64
	ReconnectTime int
	IsTest        bool
}

func (wah *WAHandler) HandleError(err error) {
	_, eMatchFailed := err.(*whatsapp.ErrConnectionFailed)
	_, eMatchClosed := err.(*whatsapp.ErrConnectionClosed)

	if eMatchFailed || eMatchClosed {
		if WASessionExist(wah.SessionFile) && wah.SessionConn != nil {
			log.Println(log.LogLevelWarn, fmt.Sprintf("connection closed unexpetedly, reconnecting after %d seconds", wah.ReconnectTime))

			<-time.After(time.Duration(wah.ReconnectTime) * time.Second)

			err := wah.SessionConn.Restore()
			if err != nil {
				log.Println(log.LogLevelError, err.Error())

				if WASessionExist(wah.SessionFile) {
					err := os.Remove(wah.SessionFile)
					if err != nil {
						log.Println(log.LogLevelError, fmt.Sprintf("remove session file error, %s", err.Error()))
					}
				}

				_, _ = wah.SessionConn.Disconnect()
				wah.SessionConn.RemoveHandlers()

				WAConn = nil
			} else {
				wah.SessionStart = uint64(time.Now().Unix())
			}
		}
	} else {
		log.Println(log.LogLevelError, err.Error())
	}
}

func (wah *WAHandler) HandleTextMessage(data whatsapp.TextMessage) {
	if wah.IsTest && data.Info.RemoteJid != wah.SessionJid {
		return
	}

	msgText := strings.SplitN(strings.TrimSpace(data.Text), " ", 2)
	if msgText[0] != wah.SessionTag || data.Info.Timestamp < wah.SessionStart {
		return
	}

	resText, err := parser.JSONExec(parser.JSONList, strings.Split(msgText[1], " "), 0)
	if err != nil {
		if len(resText) == 0 {
			resText = []string{"Ouch, Got some error here while processing your request 🙈"}
		} else {
			resText[0] = "Ouch, Got some error here while processing your request 🙈\n" + resText[0]
		}

		log.Println(log.LogLevelError, err.Error())
	}

	for i := 0; i < len(resText); i++ {
		err := WAMessageText(wah.SessionConn, data.Info.RemoteJid, resText[i], data.Info.Id, data.Text)
		if err != nil {
			log.Println(log.LogLevelError, err.Error())
		}
	}
}

func WAParseJID(jid string) string {
	components := strings.Split(jid, "@")

	if len(components) > 1 {
		jid = components[0]
	}

	suffix := "@s.whatsapp.net"

	if len(strings.SplitN(jid, "-", 2)) == 2 {
		suffix = "@g.us"
	}

	return jid + suffix
}

func WAGetSendMutexSleep() time.Duration {
	rand.Seed(time.Now().UnixNano())

	waitMin := 1000
	waitMax := 3000

	return time.Duration(rand.Intn(waitMax-rand.Intn(waitMin)) + waitMin)
}

func WASendWithMutex(conn *whatsapp.Conn, content interface{}) (string, error) {
	WACMutex.Lock()
	time.Sleep(WAGetSendMutexSleep() * time.Millisecond)

	id, err := conn.Send(content)
	WACMutex.Unlock()

	return id, err
}

func WASyncVersion(conn *whatsapp.Conn, versionClientMajor int, versionClientMinor int, versionClientBuild int) (string, error) {
	// Bug Happend When Using This Function
	// Then Set Manualy WhatsApp Client Version
	// versionServer, err := whatsapp.CheckCurrentServerVersion()
	// if err != nil {
	// 	return "", err
	// }

	conn.SetClientVersion(versionClientMajor, versionClientMinor, versionClientBuild)
	versionClient := conn.GetClientVersion()

	return fmt.Sprintf("whatsapp version %v.%v.%v", versionClient[0], versionClient[1], versionClient[2]), nil
}

func WATestPing(conn *whatsapp.Conn) error {
	ok, err := conn.AdminTest()
	if !ok {
		if err != nil {
			return err
		}

		return errors.New("something when wrong while trying to ping, please check phone connectivity")
	}

	return nil
}

func WASessionInit(versionClientMajor int, versionClientMinor int, versionClientBuild int, timeout int) (*whatsapp.Conn, string, error) {
	conn, err := whatsapp.NewConn(time.Duration(timeout) * time.Second)
	if err != nil {
		return nil, "", err
	}
	conn.SetClientName("Go WhatsApp CLI", "Go WhatsApp", "1.0")

	info, err := WASyncVersion(conn, versionClientMajor, versionClientMinor, versionClientBuild)
	if err != nil {
		return nil, "", err
	}

	WACMutex = &sync.Mutex{}

	return conn, info, nil
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
			return errors.New("session file already exist, please logout first")
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

		err = WATestPing(conn)
		if err != nil {
			return err
		}
	} else {
		return errors.New("connection is not valid")
	}

	return nil
}

func WASessionRestore(conn *whatsapp.Conn, file string) error {
	if conn != nil {
		if !WASessionExist(file) {
			_, _ = conn.Disconnect()

			return errors.New("session file doesn't exist, please login first")
		}

		session, err := WASessionLoad(file)
		if err != nil {
			_, _ = conn.Disconnect()

			if WASessionExist(file) {
				err = os.Remove(file)
				if err != nil {
					return err
				}
			}

			return errors.New("session not valid, removing session file")
		}

		session, err = conn.RestoreWithSession(session)
		if err != nil {
			_, _ = conn.Disconnect()

			return err
		}

		err = WASessionSave(file, session)
		if err != nil {
			return err
		}

		err = WATestPing(conn)
		if err != nil {
			return err
		}
	} else {
		return errors.New("connection is not valid")
	}

	return nil
}

func WASessionLogout(conn *whatsapp.Conn, file string) error {
	if conn != nil {
		defer func() {
			_, _ = conn.Disconnect()
		}()

		conn.RemoveHandlers()

		err := conn.Logout()
		if err != nil {
			return err
		}

		if WASessionExist(file) {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("connection is not valid")
	}

	return nil
}

func WASessionValidate(conn *whatsapp.Conn) error {
	if conn == nil {
		return errors.New("connection is invalid")
	}

	return nil
}

func WAMessageText(conn *whatsapp.Conn, msgJID string, msgText string, msgQuotedID string, msgQuoted string) error {
	err := WASessionValidate(conn)
	if err != nil {
		return errors.New(err.Error())
	}

	rJID := WAParseJID(msgJID)

	content := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: rJID,
		},
		Text: msgText,
	}

	if len(msgQuotedID) != 0 {
		msgQuotedProto := waproto.Message{
			Conversation: &msgQuoted,
		}

		ctxQuotedInfo := whatsapp.ContextInfo{
			QuotedMessageID: msgQuotedID,
			QuotedMessage:   &msgQuotedProto,
			Participant:     rJID,
		}

		content.ContextInfo = ctxQuotedInfo
	}

	_, err = WASendWithMutex(conn, content)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "sending message timed out":
			return nil
		default:
			return err
		}
	}

	return nil
}
