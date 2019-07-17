package mobile

import "nanomsg.org/go/mangos/v2"

func sendMessage(sock mangos.Socket, m []byte) error {
	msg := mangos.NewMessage(len(m))
	msg.Body = m
	err := sock.SendMsg(msg)
	msg.Free()
	return err
}
