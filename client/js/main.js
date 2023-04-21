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
        </div>
        `
    }
}

const index = new Index()