// class responsible for modifying the html document depending on the current route and data from the api.
class Index {
    constructor() {
        this.api = new Api()

        this.router = new Router()
        this.router.addRoute("/", this.homeHandler.bind(this))
        this.router.addRoute("/file/|fileId", this.fileHandler.bind(this))

        this.files = []
    }

    async homeHandler() {
        await this.getFiles()
        this.addContainers()
        this.addSidePanel()
        this.addHomeChatWindow()
    }

    async fileHandler(fileId) {
        await this.getFiles()

        let fileName = ""
        for (let file of this.files) {
            if (file.id == fileId) {
                fileName = file.filename
                break
            }
        }

        if (fileName == "") {
            this.router.notFoundTemplate()
            return
        }

        const previousMessages = this.api.getPreviousMessages(fileId, fileName)

        this.addContainers()
        this.addSidePanel()
        this.addFileChatWindow(fileName, previousMessages)
    }

    async getFiles() {
        try {
            this.files = await this.api.getAllFiles()
        } catch (err) {
            console.error(err)
        }
    }

    // addContainers adds the structure of html document.
    addContainers() {
        document.body.innerHTML = `
        <div id="app">
            <div id="header">
                <h1>Smartread</h1>
            </div>

            <div id="mainAppContainer">
                <div id="sidePanel"></div>
                <div id="chatWindow"></div>
            </div>
        </div>`
    }

    addSidePanel() {
        const sidePanelElm = document.getElementById("sidePanel")
        if (!sidePanelElm) {
            console.error("Side panel container must be added before running addSidePanel")
            return
        }

        // Add label for file input.
        const uploadNewFileLabelElm = document.createElement("label")
        uploadNewFileLabelElm.classList = ["fileUpload"]
        uploadNewFileLabelElm.innerText = "Upload file"
        uploadNewFileLabelElm.htmlFor = "fileUpload"
        sidePanelElm.append(uploadNewFileLabelElm)

        // Add input for uploading files.
        const uploadNewFileInputElm = document.createElement("input")
        uploadNewFileInputElm.type = "file"
        uploadNewFileInputElm.id = "fileUpload"
        uploadNewFileInputElm.accept = ".txt"
        uploadNewFileInputElm.addEventListener("change", (event) => {
            const file = event.target.files[0]
            if (!file) {
                return
            }
            this.api.uploadFile(file)
                .then(newFile => {
                    console.log("new file: ", newFile)
                    if (newFile) {
                        this.files.push(newFile)
                        this.updateFileList()
                    }
                })

        })
        sidePanelElm.append(uploadNewFileInputElm)

        // Add list of uploaded files.
        const filesContainerElm = document.createElement("div")
        filesContainerElm.id = "fileList"
        sidePanelElm.append(filesContainerElm)
        this.updateFileList()
    }

    updateFileList() {
        const filesContainerElm = document.getElementById("fileList")
        if (!filesContainerElm) {
            console.error("fileList container must be added before running updateFileList")
            return
        }
        filesContainerElm.innerHTML = ""

        for (let file of this.files) {
            const fileLinkElm = document.createElement("a")
            fileLinkElm.href = "/file/" + file.id
            fileLinkElm.innerText = file.filename
            filesContainerElm.append(fileLinkElm)
        }
    }

    addHomeChatWindow() {
        const chatWindowElm = document.getElementById("chatWindow")
        if (!chatWindowElm) {
            console.error("Chat window must be added before running addHomeChatWidow")
            return
        }

        const pElm = document.createElement("p")
        pElm.innerText = "Velg eller last opp en fil for å spørre spørsmål"
        chatWindowElm.append(pElm)
    }

    addFileChatWindow(filename, previousMessages) {
        const chatWindowElm = document.getElementById("chatWindow")
        if (!chatWindowElm) {
            console.error("Chat window must be added before running addHomeChatWidow")
            return
        }

        // Adding container for messages.
        const chatMessageContainer = document.createElement("div")
        chatMessageContainer.id = "chatMessageContainer"
        chatWindowElm.append(chatMessageContainer)

        // Adding previous messages.
        for (let message of previousMessages) {
            const messageElm = document.createElement("div")
            messageElm.classList = [message.role = "AI" ? "aiMessage" : "humanMessage", "message"]
            messageElm.innerText = message.value
            chatMessageContainer. append(messageElm)
        }

        // Adding input.
        const sendChatMessageContainer = document.createElement("div")
        sendChatMessageContainer.id = "chatInputContainer"
        const sendMessageInput = document.createElement("input")
        sendMessageInput.type = "text"
        sendChatMessageContainer.append(sendMessageInput)
        const sendMessageSubmit = document.createElement("input")
        sendMessageSubmit.type = "submit"
        
        sendChatMessageContainer.append(sendMessageSubmit)

        chatWindowElm.append(sendChatMessageContainer)
    }
}

const index = new Index()