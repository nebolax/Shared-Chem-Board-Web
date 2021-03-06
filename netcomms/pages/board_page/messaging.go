package board_page

import (
	"ChemBoard/all_boards"
)

var clients = make(map[int]sockClient)

func sendHistory(connID, boardID, viewID int) {
	if b, ok := all_boards.BoardByID(boardID); ok {
		var drawingsHist all_boards.DrawingsHistory
		var chatHist all_boards.ChatHistory
		if viewID == 0 {
			drawingsHist = b.DrawingsHistory
			chatHist = b.ChatHistory
		} else {
			if obs, ok := b.ObserverByID(viewID); ok {
				drawingsHist = obs.DrawingsHistory
				chatHist = obs.ChatHistory
			}
		}
		for _, pack := range drawingsHist {
			writeSingleMessage(connID, pack)
		}
		for _, pack := range chatHist {
			writeSingleMessage(connID, pack)
		}
	}
}

func newGroupMessage(boardID, viewID, exceptConn int, msg interface{}) {
	for _, cl := range clients {
		if cl.boardID() == boardID {
			if cl.isAdmin() {
				if cl.(adminClient).dview == viewID {
					sendtoUserDevices(cl.userID(), exceptConn, msg)
				}
			} else {
				if (!cl.(observerClient).dview && viewID == 0) || (cl.(observerClient).dview && viewID == cl.(observerClient).duserID) {
					sendtoUserDevices(cl.userID(), exceptConn, msg)
				}
			}
		}
	}
}

func sendtoUserDevices(userID, exceptConn int, message interface{}) {
	for connID, client := range clients {
		if client.userID() == userID && connID != exceptConn {
			writeSingleMessage(connID, message)
		}
	}
}

func writeSingleMessage(connID int, msg interface{}) {
	clients[connID].mu().Lock()
	defer clients[connID].mu().Unlock()

	if enc, ok := encodeMessage(msg); ok {
		err := clients[connID].sock().WriteJSON(enc)
		if err != nil {
			delClient(connID)
		}
	}
}

func delClient(connID int) {
	boardID := clients[connID].boardID()
	clients[connID].sock().Close()
	delete(clients, connID)
	updateObserversList(boardID)
}

func readSingleMessage(connID int) (interface{}, bool) {
	var msg anyMSG
	err := clients[connID].sock().ReadJSON(&msg)
	if err != nil {
		delClient(connID)
		return 0, false
	}
	return decodeMessage(msg)
}
