// Api is the class responsible for getting and sending data to the server via the api and 
// setting and getting data from local storage.
class Api {
    // Private class variables:
    #newFileEndpoint = "/api/new_file"
    #getFilesEndpoint = "/api/get_files"
    #AIRole = "AI"
    #UserRole = "User"
 
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

    getPreviousMessages(fileId, fileName) {
        let previousMessages = localStorage.getItem(fileId)
        if (previousMessages === null) {
            previousMessages = [{
                role: this.#AIRole,
                value: "Start asking me questions about the file " + fileName
            }]

            localStorage.setItem(fileId, JSON.stringify(previousMessages))
            console.log(previousMessages)
            return previousMessages
        }

        return JSON.parse(previousMessages)
    }

    saveMessage(fileId, isAIMessage, value) {
        let previousMessages = localStorage.getItem(fileId)
        previousMessages = previousMessages === null ? [] : JSON.parse(previousMessages)
       
        previousMessages.push({
            role: isAIMessage ? this.#AIRole : this.#UserRole,
            value: value,
        })

        localStorage.setItem(fileId, JSON.stringify(previousMessages))
    }

}