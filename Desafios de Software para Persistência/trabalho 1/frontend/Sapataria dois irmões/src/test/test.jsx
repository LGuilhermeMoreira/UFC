import SapatoService from '../shared/service/Sapato/SapatoService'

const Test = () => {
    return( <>
        <h1>
            Teste Service
        </h1>
        <button onClick={
           async () => {
                const a = await SapatoService.test()
                console.log(a)
            }
        }>
                test service
        </button>
        
    </>)
}

export default Test