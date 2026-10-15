import SapatoApi from "./api/SapatoApi";

class SapatoService{
    async test(){
       const response = await SapatoApi.test()
       console.log(response)
    }

    async GetSapatoById(id){
        const response = await SapatoApi.getByID(id)
        return response
    }

    async GetAllSapato(){
        return await SapatoApi.getAll()
    }

    async CreateSapato(body){
        return await SapatoApi.post(body)
    }

    async UpdateSapato(id,body){
        return await SapatoApi.put(id,body)
    }

    async DeleteSapato(id){
        return await SapatoApi.delete(id)
    }

    async GetHash(){
        return await SapatoApi.getHash()
    }

    async Download(){
        return await SapatoApi.download()
    }
}

export default new SapatoService()  

