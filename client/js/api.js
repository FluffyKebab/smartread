// Api is the class responsible for getting and sending data to the server via the api and 
// setting and getting data from local storage.
class Api {
    // Private class variables:
    #newFileEndpoint = "/api/new_file"
    #getFilesEndpoint = "/api/get_files"
    #queryEndpoint = "/api/query_file/"

    AIRole = "AI"
    userRole = "User"
 
    async uploadFile(file) {
        const formData = new FormData()
        formData.append("file", file)

        const response = await fetch(this.#newFileEndpoint, {method: "POST", body: formData,})
            .catch(err => {console.error(err)})

        if (response.status != 200) {
            alert("upload failed. Try agin later.")
            return
        }

        return await response.json()
    }

    async getAllFiles() {
        const response = await fetch(this.#getFilesEndpoint, {method: "GET"})
        if (response.status != 200) {
            throw new Error(`Error loading file data. Status: ${response.status}`)
        }
        
        const jsonResponse = await response.json()
        return jsonResponse.files
    }

    async doQuery(fileId, query) {
        const formData = new URLSearchParams()
        formData.append("query", query)

        const response = await fetch(this.#queryEndpoint + fileId, {method: "POST", body: formData})
        if (response.status != 200) {
            throw new Error(`Getting query failed: ${response.status}`)
        }
        const textResponse = await response.json()

        // Save user and AI message to local storage.
        const userMessage = {
            role: this.userRole,
            value: query,
        }

        const aiMessage = {
            role: this.AIRole,
            value: textResponse.response,
        }

        this.saveMessage(fileId, userMessage)
        this.saveMessage(fileId, aiMessage)

        return aiMessage
    }

    getPreviousMessages(fileId, fileName) {
        let previousMessages = localStorage.getItem(fileId)
        if (previousMessages === null) {
            previousMessages = [{
                role: this.AIRole,
                value: "Start asking me questions about the file " + fileName
            }]

            localStorage.setItem(fileId, JSON.stringify(previousMessages))
            console.log(previousMessages)
            return previousMessages
        }

        return JSON.parse(previousMessages)
    }

    saveMessage(fileId, message) {
        let previousMessages = localStorage.getItem(fileId)
        previousMessages = previousMessages === null ? [] : JSON.parse(previousMessages)
       
        previousMessages.push(message)
        localStorage.setItem(fileId, JSON.stringify(previousMessages))
    }

}