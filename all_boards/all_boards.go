package all_boards

import (
	"ChemBoard/netcomms/pages/account_logic"
	"ChemBoard/utils/incrementor"
	"fmt"
	"strings"
	"sync"
	"time"
)

var BoardsArray = []*DataElem{
	// {1, 1, "First", "x", []int{1, 2, 3}, sync.Mutex{}},
	// {2, 2, "Second", "y", []int{1, 2}, sync.Mutex{}},
	// {3, 1, "Third", "z", []int{2, 3}, sync.Mutex{}},
}

//TODO check if userid is valid
//TODO lock mutexes while working + boardsArray chould be private

func CreateBoard(adminID int, name, pwd string) int {
	nID := incrementor.Next("boards", true)
	board := &Board{nID, adminID, name, pwd, []*Observer{}, DrawingsHistory{}, ChatHistory{}}
	BoardsArray = append(BoardsArray, &DataElem{board, sync.Mutex{}})
	return nID
}

func BoardByID(id int) (Board, bool) {
	for _, el := range BoardsArray {
		if el.board.ID == id {
			return *el.board, true
		}
	}
	return Board{}, false
}

func (b *Board) obspointerByID(userID int) *Observer {
	for _, obs := range b.Observers {
		if obs.UserID == userID {
			return obs
		}
	}
	return nil
}

func (b Board) ObserverByID(userID int) (Observer, bool) {
	for _, obs := range b.Observers {
		if obs.UserID == userID {
			return *obs, true
		}
	}
	return Observer{}, false
}

func boardPointerByID(boardID int) *Board {
	for _, el := range BoardsArray {
		if el.board.ID == boardID {
			return el.board
		}
	}
	return nil
}

func SharedWithUser(userID int) []Board {
	res := []Board{}
	for _, el := range BoardsArray {
		for _, obs := range el.board.Observers {
			if obs.UserID == userID {
				res = append(res, *el.board)
				break
			}
		}
	}

	return res
}

func AvailableToUser(userID, boardID int) bool {
	userBoards := SharedWithUser(userID)

	if IsAdmin(userID, boardID) {
		return true
	}

	for _, b := range userBoards {
		if b.ID == boardID {
			return true
		}
	}

	return false
}

func UserAdmin(userID int) []Board {
	res := []Board{}
	for _, el := range BoardsArray {
		if el.board.Admin == userID {
			res = append(res, *el.board)
		}
	}

	return res
}

func IsAdmin(userID, boardID int) bool {
	if b, ok := BoardByID(boardID); ok && b.Admin == userID {
		return true
	}
	return false
}

func AddObserver(boardID, userID int, pwd string) bool {
	if b := boardPointerByID(boardID); b != nil {
		if b.Password == pwd {
			b.Observers = append(b.Observers, &Observer{userID, DrawingsHistory{}, ChatHistory{}})
			return true
		}
	}

	return false
}

func BoardsWithoutUser(key string, userID int) []Board {
	res := []Board{}
	for _, el := range BoardsArray {
		if strings.Contains(el.board.Name, key) && !AvailableToUser(userID, el.board.ID) {
			res = append(res, *el.board)
		}
	}

	return res
}

func NewDrawing(boardID, viewID int, msg ActionMSG) (ActionMSG, bool) {
	actionID := incrementor.Next(fmt.Sprintf("Board%d-action", boardID), true)
	drawingID := incrementor.Next(fmt.Sprintf("Board%d-drawing", boardID), true)
	msg.ID = actionID
	msg.Drawing.ID = drawingID
	bar := BoardsArray
	_ = bar
	if b := boardPointerByID(boardID); b != nil {
		if viewID == 0 {
			b.DrawingsHistory = append(b.DrawingsHistory, msg)
		} else {
			if obs := b.obspointerByID(viewID); obs != nil {
				obs.DrawingsHistory = append(obs.DrawingsHistory, msg)
			}
		}
		return msg, true
	} else {
		return ActionMSG{}, false
	}
}

func DeleteDrawing(boardID, viewID, drawingID int) {
	if b := boardPointerByID(boardID); b != nil {
		if viewID == 0 {
			res := DrawingsHistory{}
			for _, el := range b.DrawingsHistory {
				if el.ID != drawingID {
					res = append(res, el)
				}
			}
			b.DrawingsHistory = res
		} else {
			if obs := b.obspointerByID(viewID); obs != nil {
				res := DrawingsHistory{}
				for _, el := range obs.DrawingsHistory {
					if el.ID != drawingID {
						res = append(res, el)
					}
				}
				obs.DrawingsHistory = res
			}
		}
	}
}

func curTimeStamp() TimeStamp {
	ct := time.Now()
	return TimeStamp{
		ct.Year(),
		int(ct.Month()),
		ct.Day(),
		ct.Hour(),
		ct.Minute(),
	}
}

func NewChatMessage(boardID, viewID, senderID int, content ChatContent) (ChatMessage, bool) {
	if user, ok := account_logic.GetUserByID(senderID); ok {
		timeStamp := curTimeStamp()
		msgID := incrementor.Next("chat-message", true)
		msg := ChatMessage{msgID, user, timeStamp, content}

		if b := boardPointerByID(boardID); b != nil {
			if viewID == 0 {
				b.ChatHistory = append(b.ChatHistory, msg)
			} else {
				if obs := b.obspointerByID(viewID); obs != nil {
					obs.ChatHistory = append(obs.ChatHistory, msg)
				} else {
					return ChatMessage{}, false
				}
			}
		} else {
			return ChatMessage{}, false
		}

		return msg, true
	}
	return ChatMessage{}, false
}
