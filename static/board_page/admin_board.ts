module admin_board {

function initPage() {
    $("#views-nav").find("#general-page").on("click", switchView)
    $("#stepback").on("click", () => { board.stepBack() })}
    $("#exportimg").on("click", () => { board.exportPicture() })

function msgParser(e: MessageEvent) {
    let msg = JSON.parse(e.data)
    switch(msg.type) {
    case MsgTypes.Action:
        board.newAction(msg.data)
        break;
    case MsgTypes.SetId:
        switch (msg.data.property) {
            case "action":
                board.newActionID(msg.data.id)
            break;

            case "drawing":
                board.newDrawingID(msg.data.id)
            break;
        }
        break;
    case MsgTypes.ObsStat:
        msg = msg.data
        $("#observers-nav").empty()
        msg.allObsInfo.forEach((el: any) => {
            let templ = <HTMLTemplateElement>document.getElementById("template-obsname")
            let clone = document.importNode(templ.content, true)
            let btn = clone.querySelector("#chviewBtn")!!
            btn.addEventListener("click", switchView)
            btn.innerHTML = el.username
            btn.id = "view" + el.userid
            document.getElementById("observers-nav")?.appendChild(clone)
        });
        break
    case MsgTypes.InpChatMsg:
        chat.newMessage(msg.data)
        break
}
}

function switchView(e: Event) {
    let sourceId = (<HTMLElement>e.target).id
    if (sourceId == "general-page") {
        toGeneral()
    } else {
        let nview: number = +sourceId.slice(4)
        toPersonal(nview)
    }
}

function toPersonal(viewID: number) {
    board.clear()
    chat.clear()
    ws.send(JSON.stringify({
        type: MsgTypes.Chview,
        data: {
            nview: viewID
        }
    }))
}

function toGeneral() {
    board.clear()
    chat.clear()
    ws.send(JSON.stringify({
        type: MsgTypes.Chview,
        data: {
            nview: 0
        }
    }))
}

var ws: WebSocket
var board: AdminBoard
var chat: BasicChat

ws = new WebSocket('ws://' + window.location.host + "/ws" + window.location.pathname)
board = new AdminBoard(ws)
chat = new BasicChat(<HTMLDivElement>document.getElementById("chat"), ws)
initPage()
ws.onmessage = msgParser
}