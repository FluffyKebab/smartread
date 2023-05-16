// Index is the class responsible for modifying the html document depending 
// on the current route and data from the api.
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
        this.addSidePanel("")
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
        this.addSidePanel(fileId)
        this.addFileChatWindow(fileId, previousMessages)
    }

    async getFiles() {
        try {
            this.files = await this.api.getAllFiles()
        } catch (err) {
            console.error(err)
            alert("Server feil. Prøv igjen senere.")
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

    addSidePanel(fileId) {
        const sidePanelElm = document.getElementById("sidePanel")
        if (!sidePanelElm) {
            console.error("Side panel container must be added before running addSidePanel")
            return
        }

        // Add label for file input.
        const uploadNewFileLabelElm = document.createElement("label")
        uploadNewFileLabelElm.classList = ["fileUpload"]
        uploadNewFileLabelElm.innerText = "Ny fil"
        uploadNewFileLabelElm.htmlFor = "fileUpload"
        sidePanelElm.append(uploadNewFileLabelElm)

        // Add input for uploading files.
        const uploadNewFileInputElm = document.createElement("input")
        uploadNewFileInputElm.type = "file"
        uploadNewFileInputElm.id = "fileUpload"
        uploadNewFileInputElm.accept = ".txt"
        const spinner = document.createElement("img")
        spinner.src = "/img/spinner.gif"
        spinner.height = "20"
        spinner.style.visibility = "hidden"
        uploadNewFileLabelElm.append(spinner)
        uploadNewFileInputElm.addEventListener("change", (event) => {
            const file = event.target.files[0]
            if (!file) {
                return
            }
            spinner.style.visibility = "visible"

            this.api.uploadFile(file)
                .then(newFile => {
                    spinner.style.visibility = "hidden"
                    if (newFile) {
                        window.location.replace("/file/" + newFile.fileId)
                        return
                    }

                    alert("Filen kan ikke behandles. Prøv igjen senere")
                }).catch(e => {
                    console.error(e)
                    alert("Filen kan ikke behandles. Prøv igjen senere")
                })

        })
        sidePanelElm.append(uploadNewFileInputElm)

        // Legg til listen av tidligere lastet opp filer.
        const filesContainerElm = document.createElement("div")
        filesContainerElm.id = "fileList"
        sidePanelElm.append(filesContainerElm)
        this.updateFileList(fileId)
    }

    updateFileList(fileId) {
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
            if (file.id == fileId) {
                fileLinkElm.classList.add("selectedFile")
            }

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
        pElm.innerText = "Velg eller last opp en fil for å spørre spørsmål."
        chatWindowElm.append(pElm)
    }

    addFileChatWindow(fileId, previousMessages) {
        const chatWindowElm = document.getElementById("chatWindow")
        if (!chatWindowElm) {
            console.error("Chat window must be added before running addHomeChatWidow")
            return
        }

        // Lag en container for meldingene.
        const chatMessageContainer = document.createElement("div")
        chatMessageContainer.id = "chatMessageContainer"
        chatWindowElm.append(chatMessageContainer)

        // Legg til tidligere melidinger.
        for (let message of previousMessages) {
            this.addMessageToDOM(message)
        }

        // Legg inn skrive feltet.
        const chatInputBox = document.createElement("div")
        chatInputBox.id = "chatInputBox"

        const chatInputContainer = document.createElement("div")
        chatInputContainer.id = "chatInputContainer"

        const sendMessageInput = document.createElement("input")
        sendMessageInput.type = "text"
        sendMessageInput.placeholder = "Send en melding"
        chatInputContainer.append(sendMessageInput)

        const sendMessageSubmit = document.createElement("button")
        sendMessageSubmit.type = "submit"
        sendMessageInput.addEventListener("keypress", e => {
            if (e.key === "Enter") {
                e.preventDefault();
                sendMessageSubmit.click();
            }
        });
        sendMessageSubmit.onclick = e => {
            const messageValue = sendMessageInput.value
            if (messageValue == "") {
                return
            }

            sendMessageInput.value = ""
            this.addMessageToDOM({
                role: this.api.userRole,
                value: messageValue,
            })

            this.api.doQuery(fileId, messageValue)
                .then(message => this.addMessageToDOM(message))
                .catch(err => {
                    console.error(err)
                    alert("Server feil. Prøv igjen senere.")
                })
        }

        const submitIcon = document.createElement("img")
        submitIcon.src = "/img/send_message.png"
        submitIcon.height = "30"
        sendMessageSubmit.append(submitIcon)

        chatInputContainer.append(sendMessageSubmit)
        chatInputBox.append(chatInputContainer)
        chatWindowElm.append(chatInputBox)

        // Scroll ned til de nyeste meldingene.
        chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight
    }

    addMessageToDOM(message) {
        const chatMessageContainer = document.getElementById("chatMessageContainer")
        if (!chatMessageContainer) {
            console.error("chat Message Container must be added before running add message")
            return
        }

        const messageElm = document.createElement("div")
        messageElm.classList.add(
            message.role == this.api.AIRole ? "aiMessage" : "userMessage", 
            "message",
        )
        messageElm.innerText = message.value
        chatMessageContainer.append(messageElm)

        // Scroll ned til de nyeste meldingene.
        chatMessageContainer.scrollTop = chatMessageContainer.scrollHeight
    }
}

const index = new Index()