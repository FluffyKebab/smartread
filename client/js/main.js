// class responsible for modifying the html document depending on the current route and data from the api
class Index {
    constructor() {
        this.api = new Api()

        this.router = new Router()
        this.router.addRoute("/", this.homeHandler.bind(this))
        this.router.addRoute("/file/|fileId", this.fileHandler.bind(this))
    }

    homeHandler() {
        this.addContainers()
        this.addSidePanel()
    }

    fileHandler(fileId) {
        this.addContainers()
    }

    // Adds structure of html document
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

        const uploadNewFileLabelElm = document.createElement("label")
        uploadNewFileLabelElm.classList = ["fileUpload"]
        uploadNewFileLabelElm.innerText = "Upload file"
        uploadNewFileLabelElm.htmlFor = "fileUpload"
        sidePanelElm.appendChild(uploadNewFileLabelElm)

        const uploadNewFileInputElm = document.createElement("input")
        uploadNewFileInputElm.type = "file"
        uploadNewFileInputElm.id = "fileUpload"
        uploadNewFileInputElm.accept = ".txt"
        uploadNewFileInputElm.addEventListener("change", (event) => {
            const file = event.target.files[0]
            if (!file) {
                return
            }
            this.api.UploadFile(file)
        })
        sidePanelElm.appendChild(uploadNewFileInputElm)
    }
}

const index = new Index()