// Router is a simple router. 
class Router {
    constructor() {
        this.routes = {}
        this.notFoundTemplate = () => {
            document.body.innerHTML = "404 not found"
        }

        window.addEventListener('load', this.route.bind(this))
        window.addEventListener('hashchange', this.route.bind(this))
    }

    /* 
    If | is used in route characters after will be put as input to the function
    E.g /files/| /files/234 234 is input.
    */

    addRoute(path, template) {
        if (typeof template !== "function") {
            console.error("Template not function")
            return
        }

        if (typeof path !== "string") {
            console.error("Path not string")
            return
        }

        const splits = path.split("|")
        if (splits.length >= 2) {
            this.routes[splits[0] + "|"] = template
            return
        }

        this.routes[path] = template
    }

    route(evt) {
        let url = window.location.pathname

        if (this.routes.hasOwnProperty(url)) {
            this.routes[url]()
            return
        }

        for (let [path, template] of Object.entries(this.routes)) {
            if (path.endsWith("|")) {
                if (url.startsWith(path.split("|")[0])) {
                    //If current path is a open path and the url starts with the path.
                    let posOfPipe = path.indexOf("|")
                    let inputToFunction = url.slice(posOfPipe, url.length)

                    template(inputToFunction)
                    return
                }
            }
        }

        this.notFoundTemplate()
    }
}
