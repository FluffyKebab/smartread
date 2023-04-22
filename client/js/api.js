// Api is the class responsible for getting and sending data to the server via the api
class Api {
    // Private class variables:
    #newFileEndpoint = "/api/new_file"
    #getFilesEndpoint = "/api/get_files"

    uploadFile(file) {
        const formData = new FormData()
        formData.append("file", file)

        fetch(this.#newFileEndpoint, {method: "POST", body: formData,})
            .then(response => {console.log(response)})
            .catch(err => {console.error(err)})
    }

    async getAllFiles() {
        const response = await fetch(this.#getFilesEndpoint, {method: "GET"})
        if (response.status != 200) {
            throw new Error(`Error loading file data. Status: ${response.status}`)
        }
        
        const jsonResponse = await response.json()
        return jsonResponse.files
    }
}