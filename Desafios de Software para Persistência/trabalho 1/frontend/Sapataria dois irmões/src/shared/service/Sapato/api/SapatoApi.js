import { api } from "../../api/axios";

class SapatoAPI {
    async test(){
        try {
            const response = await api.get('/')
            return response
        } catch (error) {
            console.log(error)
            return error
        }
    }

    async getByID(id){
        try {
            return await api.get(`/sapato/${id}`)
        } catch (error) {
            console.log(error)
            return error
        }
    }

    async getAll(){
        try {
            return await api.get('/sapato/')
        } catch (error) {
            console.log(error)
            return error
        }
    }

    async post(body){
        try {
            return await api.post('/sapato/',body)
        } catch (error) {
            console.log(error)
            return error
        }
    }

    async put(id,body){
        try {
            return await api.put(`/sapato/${id}`,body)
        } catch (error) {
            console.log(error)
            return error
        }
    }

    async delete(id){
        try {
            return await api.delete(`/sapato/${id}`)
        } catch (error) {
            console.log(error)
            return error
        }
    }

    async getHash(){
        try {
            return await api.get(`/sapato/hash`,)
        } catch (error) {
            console.log(error)
            return error
        }
    }

    async download(){
        try {
            return await api.get(`/sapato/download-zip`,{
                responseType: 'blob'
            })
        } catch (error) {
            console.log(error)
            return error
        }
    }
}

export default new SapatoAPI()