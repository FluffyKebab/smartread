// class responsible for modifying the html document depending on the current route and data from the api
class Index {
    constructor() {
        this.api = new Api()

        this.router = new Router()
        this.router.addRoute("/", this.homeHandler.bind(this))
        this.router.addRoute("/file/|fileId", this.fileHandler.bind(this))

        this.files = []
    }

    async homeHandler() {
        await this.getData()
        this.addContainers()
        this.addSidePanel()
    }

    fileHandler(fileId) {
        this.addContainers()
    }

    async getData() {
        try {
            this.files = await this.api.getAllFiles()
        } catch (err) {
            console.error(err)
        }
    }

    // addContainers adds structure of html document
    addContainers() {
        document.body.innerHTML = `
        <div id="app">
            <div id="header">
                <h1>Smartread</h1>
            </div>

            <div id="mainAppContainer">
                <div id="sidePanel"></div>
                <div id="chat"></div>
            </div>
        </div>`
    }

    addSidePanel() {
        const sidePanelElm = document.getElementById("sidePanel")
        if (!sidePanelElm) {
            console.error("Side panel container must be added before running addSidePanel")
            return
        }

        // Add label for file input
        const uploadNewFileLabelElm = document.createElement("label")
        uploadNewFileLabelElm.classList = ["fileUpload"]
        uploadNewFileLabelElm.innerText = "Upload file"
        uploadNewFileLabelElm.htmlFor = "fileUpload"
        sidePanelElm.append(uploadNewFileLabelElm)

        // Add input for uploading files
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
        })
        sidePanelElm.append(uploadNewFileInputElm)

        // Add list of uploaded files
        const filesContainerElm = document.createElement("div")

        for (let file of this.files) {
            const fileLinkElm = document.createElement("a")
            fileLinkElm.href = "/file/" + file.id
            fileLinkElm.innerText = file.filename
            filesContainerElm.append(fileLinkElm)
        }

        sidePanelElm.append(filesContainerElm)
    }
}

const index = new Index()