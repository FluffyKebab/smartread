// Api is the class responsible for getting and sending data to the server via the api
class Api {
    // Private class variables:
    #newFileEndpoint = "/api/new_file"

    UploadFile(file) {
        const formData = new FormData()
        formData.append("file", file)

        fetch(this.#newFileEndpoint, {
            method: "POST",
            body: formData,
        }).then(response => {
            console.log(response)
        }).catch(err => {
            console.error(err)
        })
    }
}