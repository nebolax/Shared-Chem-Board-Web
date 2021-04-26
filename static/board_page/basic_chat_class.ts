class ChatMsgContent {
    text: string;

    constructor(msgContent: any) {
        this.text = msgContent.text
    }
}

class TimeStamp {
    year: number;
    month: number;
    day: number;
    hour: number;
    minute: number;

    constructor(msgStamp: any) {
        this.year = msgStamp.year
        this.month = msgStamp.month
        this.day = msgStamp.day
        this.hour = msgStamp.hour
        this.minute = msgStamp.minute
    }
}

class ChatMessage {
    id: number;
    senderid: number;
    timestamp: TimeStamp;
    content: ChatMsgContent;

    constructor(msg: any) {
        this.id = msg.id
        this.senderid = msg.senderid
        this.timestamp = new TimeStamp(msg.timestamp)
        this.content = new ChatMsgContent(msg.content)
    }
}

class BasicChat {
    history: ChatMessage[] = [];
    chatTag: HTMLDivElement;
    ws: WebSocket;
    constructor(chatTag: HTMLDivElement, ws: WebSocket) {
        this.chatTag = chatTag
        this.ws = ws
        chatTag.querySelector("#send-new-chat-msg")?.addEventListener("click", this.sendMessage)
    }

    sendMessage(e: Event) {
        let msgText = this.chatTag.querySelector("#new-chat-msg-text")?.textContent
        if (msgText == null || msgText == undefined) {
            alert("Вы должны ввести хотя бы какой-то текст")
        } else {
            this.ws.send(JSON.stringify({
                type: MsgTypes.OutChatMsg,
                data: {
                    text: msgText
                }
            }))
        }
    }
    loadHistory(msgHist: any) {
        this.history = msgHist.history
    }
    newMessage(msg: ChatMessage) {
        this.history.push(msg)
        let templ = <HTMLTemplateElement>this.chatTag.querySelector("#template-chatmsg")
        let clone = document.importNode(templ.content, true)
        clone.querySelector(".chatmsg-text")!!.innerHTML = msg.content.text
        this.chatTag.firstChild?.appendChild(clone)
    }
}